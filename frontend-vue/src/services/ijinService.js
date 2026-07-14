import api from "./api";

// GET all hari ijin
export const getAllIjin = () => api.get("/api/ijin_karyawan");

// GET presensi per tanggal
export const getIjinPerTanggal = (tanggal) => {
    return api.get("/api/ijin_karyawan/harian", {
        params: { tanggal }, // 🔥 query param
    });
};

// PUT update hari ijin berdasarkan ID
export const updateIjin = (id, data) => api.put(`/api/ijin_karyawan/${id}`, data);

// DELETE hari ijin berdasarkan ID
export const deleteIjin = (id) => api.delete(`/api/ijin_karyawan/${id}`);

// CREATE new hari ijin
export const createIjin = (data) => {
    return api.post(`/api/ijin_karyawan`,data)
}