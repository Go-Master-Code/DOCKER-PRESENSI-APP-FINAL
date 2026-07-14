import { ref } from "vue";
import * as presensiService from "@/services/presensiService";

export function usePresensi() {
  const presensiList = ref([]);
  const todayPresensi = ref(null); // catatan hari ini
  const loading = ref(false);

  // 🔥 SINGLE SOURCE OF TRUTH
  // 👉 semua filter HARUS pakai ini (tidak boleh duplikat di view)
  const selectedDate = ref(null);

  // 🔥 FETCH TANPA PARAMETER
  // 👉 kenapa? supaya tidak ada 2 sumber (param vs state)
  const fetchPresensi = async () => {
    loading.value = true;

    try {
        // 🔥 kalau ada tanggal → pakai endpoint filter
        const res = selectedDate.value
            ? await presensiService.getPresensiPerTanggal(selectedDate.value)
            : await presensiService.getAllPresensi();

        // 🔥 PENTING: pastikan selalu array
        presensiList.value = res.data.data || [];

        // 🔥 ambil data hari ini (tetap dipertahankan dari kode kamu)
        const today = new Date().toISOString().split("T")[0];

        todayPresensi.value = presensiList.value.find(p => p.tanggal === today) || null;
        } catch (err) {
            console.error(err); // 🔥 nanti bisa diganti snackbar global

            // 🔥 fallback biar UI tidak crash
            presensiList.value = [];
            todayPresensi.value = null;
        } finally {
            loading.value = false; // 🔥 wajib → hindari loading nyangkut
        }
  };

  const checkIn = async () => {
    await presensiService.createPresensi({ type: "checkin" });
    await fetchPresensi();
  };

  const checkOut = async () => {
    await presensiService.createPresensi({ type: "checkout" });
    await fetchPresensi();
  };

  // 🔥 WAJIB RETURN selectedDate
  // 👉 supaya view & composable pakai state yang sama
  return { presensiList, todayPresensi, loading, selectedDate, fetchPresensi, checkIn, checkOut };
}