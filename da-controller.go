package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// @Summary Returns Celestia DA node info
// @Description Pings Celestia DA node and returns node info
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /node_info [get]
func NodeInfoController(c *fiber.Ctx) error {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	l.Println("celestia da rpc node ping...")

	// get rpc da info
	headerInfo, err := NodePing(ctx, celestiaRpcAddress, authToken)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	response := Response{
		Success: true,
		Result: fiber.Map{
			"height":     headerInfo.Height(),
			"chain_id":   headerInfo.ChainID(),
			"block_hash": headerInfo.Hash(),
			"timestamp":  headerInfo.Time(),
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func SubmitJSONDataController(c *fiber.Ctx) error {
	var payload Data

	// 요청 본문을 JSON으로 파싱
	if err := c.BodyParser(&payload); err != nil {
		response := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if payload.NamespaceKey == "" {
		response := Response{
			Success: false,
			Error:   "namespace key is required for celestia da",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	encodedNamesapceKey := base64.StdEncoding.EncodeToString([]byte(payload.NamespaceKey))
	namespaceKey := []byte(encodedNamesapceKey)[:10]

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// submit data
	height, err := SubmitBlob(ctx, celestiaRpcAddress, authToken, namespaceKey, payload.MetaData)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	l.Infof("Metadata Blob was included at height %d in %X namespace\n", height, namespaceKey)

	// retrieve blobs
	retrievedBlobs, err := GetBlobs(ctx, celestiaRpcAddress, authToken, height)
	if err != nil {
		resp := Response{
			Success: false,
			Error:   fmt.Sprintf("Successfully, submitted formdata to celestia da but, failed to retrieve blobs in the height. unexpected err: %s", err.Error()),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	var hashStr string
	for _, blob := range retrievedBlobs {
		if bytes.Equal(blob.Data, payload.MetaData) {
			l.Infof("found matched metadata in blobs in %d height!", height)
			l.Printf("blob commitment hash: %X \n", blob.Commitment)
			// l.Printf("blob commitment: %s \n", blob.Commitment)
			hashStr = fmt.Sprintf("%X", blob.Commitment)
			break
		}
	}

	// insert db hash for saving height
	index, err := InsertNamespace(db, string(namespaceKey), hashStr, int(height))
	if err != nil {
		resp := Response{
			Success: false,
			Error:   fmt.Sprintf("Successfully, submitted formdata to celestia da but, failed to save commitment hash & height in db. unexpected err: %s", err.Error()),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	response := Response{
		Success: true,
		Result: fiber.Map{
			"namespace_key":            string(namespaceKey),
			"submitted_metadata_index": index,
			"submitted_metadata":       payload.Data,
			"submitted_height":         height,
		},
	}
	return c.JSON(response)
}

func SubmitFormDataController(c *fiber.Ctx) error {
	// 파일 업로드 처리
	file, err := c.FormFile("image")
	if err != nil {
		l.Errorln(err)

		response := Response{
			Success: false,
			Error:   "cannot parse form file",
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// 파일 읽기
	f, err := file.Open()
	if err != nil {
		response := Response{
			Success: false,
			Error:   "cannot open uploaded file",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		response := Response{
			Success: false,
			Error:   "cannot read uploaded file",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// 기본 정보를 로그로 출력
	l.Infof("received file: %s", file.Filename)
	l.Infof("file size: %d bytes", len(fileBytes))

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// submit data
	height, err := SubmitBlobImage(ctx, celestiaRpcAddress, authToken, fileBytes)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	l.Infof("FormData Blob was included at height %d\n", height)

	// retrieve blobs
	retrievedBlobs, err := GetBlobs(ctx, celestiaRpcAddress, authToken, height)
	if err != nil {
		// l.Errorf("Successfully, submitted formdata to celestia da but, failed to retrieve blobs in the height. unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   fmt.Sprintf("Successfully, submitted formdata to celestia da but, failed to retrieve blobs in the height. unexpected err: %s", err.Error()),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	var hashStr string
	for _, blob := range retrievedBlobs {
		if bytes.Equal(blob.Data, fileBytes) {
			// l.Infoln("equal!")
			l.Printf("blob commitment hash: %X \n", blob.Commitment)
			// l.Printf("blob commitment: %s \n", blob.Commitment)
			hashStr = fmt.Sprintf("%X", blob.Commitment)
			break
		}
	}

	// insert db hash for saving height
	err = InsertBlob(db, hashStr, int(height))
	if err != nil {
		resp := Response{
			Success: false,
			Error:   fmt.Sprintf("Successfully, submitted formdata to celestia da but, failed to save commitment hash & height in db. unexpected err: %s", err.Error()),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	response := Response{
		Success: true,
		Result: fiber.Map{
			"hash":             hashStr,
			"submitted_height": height,
		},
	}

	return c.JSON(response)
}

func RetrieveBlobController(c *fiber.Ctx) error {
	retrieveHeightStr := c.Query("height")

	l.Infof("received height: %s", retrieveHeightStr)

	retrieveHeight, err := strconv.ParseUint(retrieveHeightStr, 10, 64)
	if err != nil {
		response := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// submit data
	retrievedBlobs, err := GetBlobs(ctx, celestiaRpcAddress, authToken, retrieveHeight)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	var data []byte

	// 각 블롭을 직렬화하고 리스트에 추가 (그냥 하나만 있다고 가정)
	for _, blob := range retrievedBlobs {
		l.Printf("blob commitment: %X \n", blob.Commitment)
		l.Printf("blob Namespace: %X \n", blob.Namespace)
		l.Printf("blob NamespaceVersion: %d \n", blob.NamespaceVersion)
		l.Printf("blob Data: %d \n", len(blob.Data))
		l.Printf("blob index: %d \n", blob.Index)
		l.Printf("blob ShareVersion: %d \n", blob.ShareVersion)

		data = blob.Data
		break
	}

	return c.Send(data)
}

func RetrieveBlobByCommitment(c *fiber.Ctx) error {
	hashStr := c.Params("hash") // URL 파라미터에서 해시 값 가져오기
	l.Infof("Received hash: %s\n", hashStr)

	retrieveHeight, err := GetBlobHeight(db, hashStr)
	if err != nil {
		return err
	}

	l.Infof("Found height by hash: %d\n", retrieveHeight)

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// submit data
	blob, err := GetBlob(ctx, celestiaRpcAddress, authToken, uint64(retrieveHeight), hashStr)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	l.Printf("blob commitment: %X \n", blob.Commitment)
	l.Printf("blob Namespace: %X \n", blob.Namespace)
	l.Printf("blob NamespaceVersion: %d \n", blob.NamespaceVersion)
	l.Printf("blob Data: %d \n", len(blob.Data))
	l.Printf("blob index: %d \n", blob.Index)

	return c.Send(blob.Data)
}

func RetrieveBlobByNamespaceKey(c *fiber.Ctx) error {
	namespace := c.Params("namespace")
	index := c.Params("index")

	l.Infof(namespace, index)

	var namespaceIndex int
	_, err := fmt.Sscanf(index, "%d", &namespaceIndex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid index",
		})
	}

	hashStr, height, err := GetNamespaceData(db, namespace, namespaceIndex)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "data not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// get blob by height and commitment hash
	blob, err := GetBlob(ctx, celestiaRpcAddress, authToken, uint64(height), hashStr)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	l.Printf("blob commitment: %X \n", blob.Commitment)
	l.Printf("blob Namespace: %X \n", blob.Namespace)
	l.Printf("blob NamespaceVersion: %d \n", blob.NamespaceVersion)
	l.Printf("blob Data: %d \n", len(blob.Data))
	l.Printf("blob Data: %s \n", (blob.Data))
	l.Printf("blob index: %d \n", blob.Index)

	return c.Send(blob.Data)
}