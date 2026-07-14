import api from "./api";

// GET semua presensi
//export const getAllPresensi = () => api.get("/api/presensi");

// GET presensi hari ini
export const getAllPresensi = () => api.get("/api/presensi/today");

// GET presensi per tanggal
export const getPresensiPerTanggal = (tanggal) => {
    return api.get("/api/presensi/harian", {
        params: { tanggal }, // 🔥 query param
    });
};

// POST presensi (check-in / check-out), berfungsi juga untuk update waktu pulang secara manual di dashboard
export const createPresensi = (payload) => api.post("/api/presensi", payload);

// endpoint untuk tampilkan report presensi all karyawan per periode tanggal (jml hadir dan tidak hadir per row karyawan)
export const getLaporanPeriode = (params) => {
    return api.get("/api/presensi/all/periode", {
        params : {
            awal: params.awal,
            akhir: params.akhir,
        },
    });
}

// endpoint untuk tampilkan report presensi per karyawan per periode tanggal
export const getLaporanKaryawanPerPeriode = (params) => {
    return api.get("/api/presensi/karyawan/periode", {
        params : {
            id: params.id,
            awal: params.awal,
            akhir: params.akhir,
        }
    })
}

// endpoint get jumlah hari kerja untuk report presensi per periode
export const getHariKerja = (params) => {
    return api.get("/api/hari_kerja", {
        params : {
            awal: params.awal,
            akhir: params.akhir,
        },
    });
}