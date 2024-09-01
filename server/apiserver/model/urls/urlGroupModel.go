package urls

import (
	"context"
	"fmt"

	"github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/model/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func init() {
	urlFields := "id, group_id, title, url, created_at"

	urlTemplates := map[string]string{
		"add_urls_v1":               "select _id from main.add_urls_v1($1, $2, $3, $4)",
		"get_distinct_groups_count": "select count(id), group_id from main.urls group by group_id",
		"get_urls_v1_by_user_id":    fmt.Sprintf("select %s from main.urls where user_id = $1", urlFields),
	}

	database.NewStatements().AddSchemaTemplateMap(urlTemplates)
}

func (urlModel *UrlModel) create(ctx context.Context, urls *Urls, tx pgx.Tx) error {
	createUrlsSql := database.NewStatements().GetSchemaTemplate("add_urls_v1")
	if createUrlsSql == "" {
		return fmt.Errorf("add_urls_v1 SQL function not found")
	}

	batch := &pgx.Batch{}
	for _, item := range urls.UrlRequest {
		batch.Queue(createUrlsSql, item.GroupId, item.UserId, item.Url, item.Title)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < batch.Len(); i++ {
		var id string
		err := results.QueryRow().Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to insert URL: %w", err)
		}
	}

	return nil
}

func (urls *Urls) CreateUrls(ctx *gin.Context) (string, error) {
	urlModel := &UrlModel{}

	err := ctx.ShouldBindJSON(&urls.UrlRequest)
	if err != nil {
		return "", controller.ErrInvalidInput
	}

	groupId := uuid.New().String()

	userID := urls.Auth.UserId
	if userID == "" {
		return "", controller.ErrUserNotAuthorized
	}

	for i := range urls.UrlRequest {
		urls.UrlRequest[i].UserId = userID
		urls.UrlRequest[i].GroupId = groupId
	}

	tx, err := urls.Database.CPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
		return "", controller.ErrInternal
	}

	defer func() {
		if err != nil {
			txErr := tx.Rollback(ctx)
			if txErr != nil {
				log.Error().Err(txErr).Msg("failed to rollback transaction")
			}
		} else {
			txErr := tx.Commit(ctx)
			if txErr != nil {
				log.Error().Err(txErr).Msg("failed to commit transaction")
				err = txErr
			} else {
				log.Info().Msgf("successfully created urls")
			}
		}
	}()

	err = urlModel.create(ctx, urls, tx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create urls")
		return "", err
	}

	return groupId, nil
}

func (urlModel *UrlModel) getUrlsByUserId(ctx context.Context, urls *Urls, tx pgx.Tx) error {
	getUrlsByUserIdSql := database.NewStatements().GetSchemaTemplate("get_urls_v1_by_user_id")
	if getUrlsByUserIdSql == "" {
		return fmt.Errorf("get_urls_v1_by_user_id SQL function not found")
	}

	rows, err := tx.Query(ctx, getUrlsByUserIdSql, urls.Auth.UserId)
	if err != nil {
		return fmt.Errorf("failed to fetch URLs: %w", err)
	}

	defer rows.Close()

	var urlResponse []UrlModel
	for rows.Next() {
		var urlResponseItem UrlModel
		err := rows.Scan(&urlResponseItem.Id, &urlResponseItem.GroupId, &urlResponseItem.Title, &urlResponseItem.Url, &urlResponseItem.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan URL row: %w", err)
		}

		urlResponse = append(urlResponse, urlResponseItem)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over URL rows: %w", err)
	}

	urls.UrlResponse = urlResponse
	return nil
}

func (urls *Urls) GetUrlsByUserId(ctx *gin.Context) (*[]TransformUrls, error) {
	urlModel := &UrlModel{}

	if urls.Auth.UserId == "" {
		return nil, controller.ErrUserNotAuthorized
	}

	tx, err := urls.Database.CPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
		return nil, controller.ErrInternal
	}

	defer func() {
		if err != nil {
			txErr := tx.Rollback(ctx)
			if txErr != nil {
				log.Error().Err(txErr).Msg("failed to rollback transaction")
			}
		} else {
			txErr := tx.Commit(ctx)
			if txErr != nil {
				log.Error().Err(txErr).Msg("failed to commit transaction")
				err = txErr
			} else {
				log.Info().Msgf("successfully fetched urls")
			}
		}
	}()

	err = urlModel.getUrlsByUserId(ctx, urls, tx)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch urls")
		return nil, err
	}

	currentGroupId := ""

	var tabs []Tabs
	var transformedUrls []TransformUrls

	if urls.UrlResponse != nil {
		for index, item := range urls.UrlResponse {
			if currentGroupId == "" {
				currentGroupId = item.GroupId
			}

			if currentGroupId == item.GroupId {
				tab := Tabs{
					Id:    item.Id,
					Url:   item.Url,
					Title: item.Title,
				}

				tabs = append(tabs, tab)

				if index+1 >= len(urls.UrlResponse) {
					transformedUrls = append(transformedUrls, TransformUrls{
						GroupId: currentGroupId,
						Tabs:    tabs,
					})
				} else {
					if urls.UrlResponse[index+1].GroupId != currentGroupId {
						transformedUrls = append(transformedUrls, TransformUrls{
							GroupId: currentGroupId,
							Tabs:    tabs,
						})

						currentGroupId = ""
						tabs = nil
					}
				}

			}

		}

		return &transformedUrls, nil
	}

	return nil, err

}
