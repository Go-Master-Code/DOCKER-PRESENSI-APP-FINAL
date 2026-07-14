import api from "./api";

// GET all karyawan
export const getAllKaryawan = () => api.get("/api/karyawan");

// GET hanya karyawan yang belum absen di tanggal tertentu
export const getAllKaryawanBelumAbsen = (tanggal) => api.get(`/api/karyawan/absen/${tanggal}`);

// GET hanya karyawan yang belum ijin di tangga tertentu
export const getAllKaryawanBelumIjin = (tanggal) => api.get(`/api/karyawan/ijin/${tanggal}`);

// PUT update karyawan berdasarkan ID
export const updateKaryawan = (id, data) => api.put(`/api/karyawan/${id}`, data);

// DELETE karyawan berdasarkan ID
export const deleteKaryawan = (id) => api.delete(`/api/karyawan/${id}`);

// CREATE new karyawan
export const createKaryawan = (data) => {
    return api.post(`/api/karyawan`,data)
}

// IMPORT data from excel
export const importExcel = (formData) => {
  return api.post("api/karyawan/import", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
};