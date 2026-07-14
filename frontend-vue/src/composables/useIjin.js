import { ref } from "vue";
import * as ijinService from "@/services/ijinService";

export function useIjin() {
  const ijinList = ref([]);
  const todayIjin = ref(null); // catatan hari ini
  const loading = ref(false);

  // 🔥 SINGLE SOURCE OF TRUTH
  // 👉 semua filter HARUS pakai ini (tidak boleh duplikat di view)
  const selectedDate = ref(null);

  // 🔥 FETCH TANPA PARAMETER
  // 👉 kenapa? supaya tidak ada 2 sumber (param vs state)
  const fetchIjin = async () => {
      loading.value = true;

      try {
          // 🔥 GUNAKAN selectedDate LANGSUNG
          // 👉 ini inti dari fix: filter selalu konsisten
          const res = selectedDate.value
              ? await ijinService.getIjinPerTanggal(selectedDate.value)
              : await ijinService.getAllIjin();

          // 🔥 PENTING: pastikan selalu array
          ijinList.value = res.data.data || [];

          // 🔥 ambil data hari ini (tetap dipertahankan dari kode kamu)
          const today = new Date().toISOString().split("T")[0];

          todayIjin.value = ijinList.value.find(i => i.tanggal === today) || null;
          } catch (err) {
              console.error(err); // 🔥 nanti bisa diganti snackbar global

              // 🔥 fallback biar UI tidak crash
              ijinList.value = [];
              todayIjin.value = null;
          } finally {
              loading.value = false; // 🔥 wajib → hindari loading nyangkut
          }
    };
  
  return { ijinList, todayIjin, loading, fetchIjin, selectedDate };
  // 🔥 WAJIB RETURN selectedDate
  // 👉 supaya view & composable pakai state yang sama
}