import api from "./api";

// GET all user
export const getAllUser = () => api.get("/api/user");

// PUT update user berdasarkan ID
export const updateUser = (id, data) => api.put(`/api/user/${id}`, data);

// DELETE user berdasarkan ID
export const deleteUser = (id) => api.delete(`/api/user/${id}`);

// CREATE new user
export const createUser = (data) => {
    return api.post(`/api/user`,data)
}