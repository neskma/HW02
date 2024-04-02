package v1

import (
	"net/http"

	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/api/apiv1"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/conv"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/httputil"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/pb"
)

func newUsersHandler(usersClient usersClient) *usersHandler {
	return &usersHandler{client: usersClient}
}

type usersHandler struct {
	client usersClient
}

func (h *usersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp, err := h.client.ListUsers(ctx, &pb.Empty{})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Преобразование списка пользователей в API v1 формат
	userList := make([]apiv1.User, 0, len(resp.Users))
	for _, u := range resp.Users {
		userList = append(userList, conv.ToAPIV1User(u))
	}

	// Отправка ответа с данными пользователей
	httputil.MarshalResponse(w, http.StatusOK, userList)
}

func (h *usersHandler) PostUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u apiv1.UserCreate
	code, err := httputil.Unmarshal(w, r, &u)
	if err != nil {
		httputil.MarshalResponse(w, code, apiv1.Error{
			Code:    httputil.ConvertHTTPToErrorCode(code),
			Message: conv.ToPtr(err.Error()),
		})
		return
	}

	// Создание нового пользователя через gRPC клиент
	if _, err := h.client.CreateUser(ctx, conv.ToGRPCUserCreate(u)); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *usersHandler) DeleteUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	// Удаление пользователя через gRPC клиент
	if _, err := h.client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id}); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *usersHandler) GetUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	// Получение пользователя по ID через gRPC клиент
	u, err := h.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Преобразование пользователя в API v1 формат
	user := conv.ToAPIV1User(*u)

	// Отправка ответа с данными пользователя
	httputil.MarshalResponse(w, http.StatusOK, user)
}

func (h *usersHandler) PutUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var u apiv1.UserCreate
	code, err := httputil.Unmarshal(w, r, &u)
	if err != nil {
		httputil.MarshalResponse(w, code, apiv1.Error{
			Code:    httputil.ConvertHTTPToErrorCode(code),
			Message: conv.ToPtr(err.Error()),
		})
		return
	}

	// Обновление информации о пользователе через gRPC клиент
	if _, err := h.client.UpdateUser(ctx, conv.ToGRPCUserUpdate(id, u)); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleGRPCError(w http.ResponseWriter, err error) {
	// Обработка ошибки gRPC
	httputil.MarshalResponse(w, http.StatusInternalServerError, apiv1.Error{
		Code:    http.StatusInternalServerError,
		Message: conv.ToPtr(err.Error()),
	})
}
