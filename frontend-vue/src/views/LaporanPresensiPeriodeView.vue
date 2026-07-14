<script setup>
    import {ref, watch, onMounted} from "vue";
    import dayjs from "dayjs";
    import * as presensiService from "@/services/presensiService";
    // import library excel -> npm install xlsx file-saver di console
    import * as XLSX from "xlsx";
    import { saveAs } from "file-saver";

    // state filter
    const awal = ref("");
    const akhir = ref("");
    // state jumlah hari kerja
    const hariKerja = ref(null);
    
    // ==============================
    // 🔥 ERROR STATE (INLINE)
    // ==============================
    const errors = ref({
      tanggal: "",
      fetch: ""
    });

    // ==============================
    // 🔥 FUNCTION untuk HITUNG jumlah TIDAK HADIR (FRONTEND)
    // ==============================
    const getTidakHadir = (item) => {
      const hadir = item.kehadiran || 0;
      const ijin = item.jumlah_ijin || 0;
      const totalHari = hariKerja.value || 0;

      // 🔥 SAFE GUARD → tidak boleh negatif
      return Math.max(0, totalHari - hadir - ijin);
    };

    // watcher tanggal awal dan akhir, auto-clear error saat user input
    watch([awal, akhir], () => {
      errors.value.tanggal = "";
    });

    // default value (UX bagus)
    onMounted(() => {
        awal.value = dayjs().startOf("month").format("YYYY-MM-DD"); // tglAwal diset ke tgl awal bulan
        akhir.value = dayjs().format("YYYY-MM-DD"); // tglAkhir diset ke tgl hari ini
    }); 

    // state data
    const loading = ref(false);
    const items = ref([]);     // 🔥 untuk tabel (ARRAY)

    // table header
    const headers = [
        { title:"ID", key: "karyawan_id" },
        { title:"Nama", key: "nama" },
        { title:"Jumlah Hadir", key: "kehadiran" },
        { title:"Jumlah Ijin", key: "jumlah_ijin" },
        { title:"Tidak Hadir", key: "tidak_hadir" },
    ]

    // === PAGINATION ===
    // 🔥 v-data-table options (controlled pagination)
    const options = ref({
        page: 1,
        itemsPerPage: 5,
    });

    // fetch laporan
    const fetchLaporan = async() => {
        // 🔥 reset error
        errors.value = {
          tanggal: "",
          fetch: ""
        };

        // validasi input input awal dan akhir
        if (!awal.value || !akhir.value) {
          errors.value.tanggal = "Tanggal awal dan akhir wajib diisi";
          return;
        }

        // validasi range tanggal
        if (awal.value > akhir.value) {
            errors.value.tanggal = "Tanggal awal tidak boleh lebih besar dari tanggal akhir";
            return;
        }

        // mulai try and catch
        try {
            loading.value = true;

            // call api paralel (2 sekaligus yaitu getHariKerja dan getLaporanPerPeriode)
            const [laporanRes, hariKerjaRes] = await Promise.all([
                // urutan api / endpoint harus sesuai urutan dengan const [] di atas!
                // record data dulu, baru hariKerja
                // api laporan per periode per karyawan
                presensiService.getLaporanPeriode({
                    awal : awal.value,
                    akhir : akhir.value,
                }),

                // api hari kerja
                presensiService.getHariKerja({
                    awal : awal.value,
                    akhir : akhir.value,
                }),
            ]);

            // ambil data dari response json
            items.value = laporanRes.data.data || []; // jika data null, maka kirim map kosong
            // set data hari kerja
            hariKerja.value = hariKerjaRes.data.data || 0; // default value hari kerja = 0
            // DEBUG row data dan jumlah hari kerja
            // console.log("LAPORAN:", laporanRes.data.data);
            // console.log("TYPE:", typeof laporanRes.data.data);
          } catch (err) {
            console.error(err);
            items.value = [];
            hariKerja.value = null;
            // 🔥 tampilkan error API
            errors.value.fetch = err.response?.data?.error || "Gagal mengambil data laporan";
        } finally {
            loading.value = false;
        }
    };

    // ==============================
    // 🔥 CLEAR FORM (RESET KE DEFAULT)
    // ==============================
    const clearForm = () => {
      // 🔄 reset tanggal ke default (sama seperti onMounted)
      awal.value = dayjs().startOf("month").format("YYYY-MM-DD"); // tanggal awal di bulan ini
      akhir.value = dayjs().format("YYYY-MM-DD"); // tanggal hari ini

      // 🔄 reset data tabel
      items.value = [];

       // 🔥 reset hari kerja → ini kunci
      hariKerja.value = null;

      // 🔄 reset error
      errors.value = {
        tanggal: "",
        fetch: "",
      };
    };

    const exportExcel = () => {
      // 🔥 guard: kalau belum ada data
      if (!items.value.length) {
        errors.value.fetch = "Tidak ada data untuk diexport";
        return;
      }

      // 🔥 mapping data → format Excel
      const data = items.value.map(item => ({
        "ID Karyawan": item.karyawan_id,
        "Nama": item.nama,
        "Jumlah Hadir": item.kehadiran,
        "Jumlah Ijin": item.jumlah_ijin,
        "Tidak Hadir": getTidakHadir(item),
      }));

      // 🔥 buat worksheet
      const worksheet = XLSX.utils.json_to_sheet(data, { origin: "A4" });

      // baru tambahkan header di atas
      XLSX.utils.sheet_add_aoa(
        worksheet,
        [
          ["Laporan Presensi Karyawan"],
          [`Periode: ${awal.value} s/d ${akhir.value}`],
          []
        ],
        { origin: "A1" }
      );

      // auto width column
      worksheet["!cols"] = [
        { wch: 12 }, // ID
        { wch: 25 }, // Nama
        { wch: 12 }, // Jumlah hadir
        { wch: 12 }, // Jumlah ijin
        { wch: 12 }, // Jumlah tidak hadir
      ];

      // 🔥 buat workbook
      const workbook = XLSX.utils.book_new();
      XLSX.utils.book_append_sheet(workbook, worksheet, "Laporan");

      // 🔥 nama file dinamis
      const fileName = `Laporan_Presensi_All_${awal.value}_sd_${akhir.value}.xlsx`;

      // 🔥 generate file
      const excelBuffer = XLSX.write(workbook, {
        bookType: "xlsx",
        type: "array",
      });

      const blob = new Blob([excelBuffer], {
        type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8",
      });

      saveAs(blob, fileName);
    };

</script>

<template>
  <v-container fluid>

    <!-- ========================= -->
    <!-- 🔥 FILTER SECTION -->
    <!-- ========================= -->
    <v-card class="pa-4 mb-4" elevation="4" rounded="xl">
      <v-card-title class="d-flex flex-column gap-2">
        <div class="text-h3 font-weight-bold d-flex gap-2">
            <v-icon size="28">mdi-calendar-range</v-icon>
            Laporan Presensi Per Periode
        </div>
      </v-card-title>

      <div class="d-flex align-start gap-3 flex-wrap" style="gap: 12px;">
          <!-- simple & clean (native date input) -->
          <v-text-field
              v-model="awal"
              label="Tanggal Awal"
              type="date"
              density="compact"
              variant="outlined"
              style="max-width: 180px"
          />

          <!-- 🔥 TANGGAL AKHIR -->
            <v-text-field
              v-model="akhir"
              label="Tanggal Akhir"
              type="date"
              density="compact"
              variant="outlined"
              style="max-width: 180px"

              :error="!!errors.tanggal"
              :error-messages="errors.tanggal"
            />

          <!-- 🔥 BUTTON FILTER -->
          <v-btn
              color="primary"
              :loading="loading"
              :disabled="!awal || !akhir"
              @click="fetchLaporan"
              height="40">
              <v-icon start>mdi-filter</v-icon>
              Filter
          </v-btn>

          <!-- 🔥 BUTTON CLEAR -->
          <v-btn
              color="warning"
              :loading="loading"
              @click="clearForm"
              height="40">
              <v-icon start>mdi-refresh</v-icon>
              Clear
          </v-btn>

          <!-- 🔥 BUTTON EXPORT TO EXCEL -->
          <v-btn
            color="success"
            :disabled="!items.length"
            @click="exportExcel"
            height="40"
          >
            <v-icon start>mdi-file-excel</v-icon>
            Export
          </v-btn>

      </div>

    </v-card>

    <!-- ========================= -->
    <!-- ❌ ERROR FETCH -->
    <!-- ========================= -->
    <div v-if="errors.fetch" class="mb-3">
      <v-alert
        type="error"
        density="compact"
        variant="tonal"
        closable
        @click:close="errors.fetch = ''"
      >
        {{ errors.fetch }}
      </v-alert>
    </div>

    <!-- ========================= -->
    <!-- 🔥 INFO HARI KERJA -->
    <!-- ========================= -->
    <div v-if="hariKerja !== null" class="mb-3">
      <v-alert
        type="info"
        density="compact"
        variant="tonal"
      >
        📅 Jumlah hari kerja:
        <strong>{{ hariKerja }}</strong> hari
      </v-alert>

    </div>

    <!-- ========================= -->
    <!-- 🔥 TABLE -->
    <!-- ========================= -->
    <v-card elevation="6" rounded="xl">

      <v-data-table
        v-model:options="options"
        :headers="headers"
        :items="items"
        :loading="loading"
        :items-per-page-options="[5, 10, 25, 50, { title: 'All', value: -1 }]"
        class="modern-table"
        hover
        density="comfortable"
        rounded="lg"
      >

        <!-- 🔥 LOADING -->
        <template #loading>
          <v-skeleton-loader type="table-row@5" />
        </template>

        <!-- JUMLAH ALPHA / TIDAK HADIR -->
        <template #item.tidak_hadir="{ item }">
          <v-chip
            :color="getTidakHadir(item) > 0 ? 'red' : 'green'"
            size="small"
            variant="tonal"
          >
            {{ getTidakHadir(item) }}
          </v-chip>
        </template>

        <!-- 🔥 EMPTY STATE -->
        <template #no-data>
          <div class="text-center pa-4 text-grey">
            Tidak ada data
          </div>
        </template>

      </v-data-table>

    </v-card>

  </v-container>
</template>