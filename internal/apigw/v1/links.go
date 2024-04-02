package v1

import (
	"net/http"

	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/api/apiv1"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/conv"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/httputil"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/pb"
)

func newLinksHandler(linksClient linksClient) *linksHandler {
	return &linksHandler{client: linksClient}
}

type linksHandler struct {
	client linksClient
}

func (h *linksHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp, err := h.client.ListLinks(ctx, nil)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Преобразование списка ссылок в API v1 формат
	links := conv.ToAPIV1Links(resp.Links)

	// Отправка ответа с данными ссылок
	httputil.MarshalResponse(w, http.StatusOK, links)
}

func (h *linksHandler) PostLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var link apiv1.LinkCreate
	code, err := httputil.Unmarshal(w, r, &link)
	if err != nil {
		httputil.MarshalResponse(
			w, code, apiv1.Error{
				Code:    httputil.ConvertHTTPToErrorCode(code),
				Message: conv.ToPtr(err.Error()),
			},
		)
		return
	}

	// Создание новой ссылки через gRPC клиент
	if _, err := h.client.CreateLink(ctx, conv.ToGRPCCreateLinkRequest(link)); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *linksHandler) DeleteLinksId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	// Удаление ссылки по ID через gRPC клиент
	if _, err := h.client.DeleteLink(ctx, &pb.DeleteLinkRequest{Id: id}); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *linksHandler) GetLinksId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	// Получение ссылки по ID через gRPC клиент
	link, err := h.client.GetLink(ctx, &pb.GetLinkRequest{Id: id})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Преобразование ссылки в API v1 формат
	apiLink := conv.ToAPIV1Link(*link)

	// Отправка ответа с данными ссылки
	httputil.MarshalResponse(w, http.StatusOK, apiLink)
}

func (h *linksHandler) PutLinksId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var link apiv1.LinkCreate
	code, err := httputil.Unmarshal(w, r, &link)
	if err != nil {
		httputil.MarshalResponse(
			w, code, apiv1.Error{
				Code:    httputil.ConvertHTTPToErrorCode(code),
				Message: conv.ToPtr(err.Error()),
			},
		)
		return
	}

	// Обновление информации о ссылке через gRPC клиент
	if _, err := h.client.UpdateLink(ctx, conv.ToGRPCUpdateLinkRequest(id, link)); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *linksHandler) GetLinksUserUserID(w http.ResponseWriter, r *http.Request, userID string) {
	ctx := r.Context()

	// Получение списка ссылок по ID пользователя через gRPC клиент
	resp, err := h.client.GetLinksByUserID(ctx, &pb.GetLinksByUserId{UserId: userID})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Преобразование списка ссылок в API v1 формат
	links := conv.ToAPIV1Links(resp.Links)

	// Отправка ответа с данными ссылок
	httputil.MarshalResponse(w, http.StatusOK, links)
}
