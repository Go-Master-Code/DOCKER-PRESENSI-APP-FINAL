import api from "./api";

// GET all jenis ijin
export const getAllJenisIjin = () => api.get("/api/jenis_ijin");

// PUT update jenis ijin berdasarkan ID
export const updateJenisIjin = (id, data) => api.put(`/api/jenis_ijin/${id}`, data);

// DELETE jenis ijin berdasarkan ID
export const deleteJenisIjin = (id) => api.delete(`/api/jenis_ijin/${id}`);

// CREATE new jenis_ijin
export const createJenisIjin = (data) => {
    return api.post(`/api/jenis_ijin`,data)
}