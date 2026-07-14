import api from "./api";

// GET all hari libur
export const getAllHariLibur = () => api.get("/api/hari_libur");

// PUT update hari libur berdasarkan ID
export const updateHariLibur = (id, data) => api.put(`/api/hari_libur/${id}`, data);

// DELETE hari libur berdasarkan ID
export const deleteHariLibur = (id) => api.delete(`/api/hari_libur/${id}`);

// CREATE new hari libur
export const createHariLibur = (data) => {
    return api.post(`/api/hari_libur`,data)
}