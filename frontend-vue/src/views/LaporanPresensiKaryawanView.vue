<script setup>
    import {ref, watch, computed, onMounted} from "vue";
    import dayjs from "dayjs";
    import * as presensiService from "@/services/presensiService";
    // import service karyawan untuk get data dan dimasukkan ke dalam combobox
    import * as karyawanService from "@/services/karyawanService";
    // import library excel -> npm install xlsx file-saver di console
    import * as XLSX from "xlsx";
    import { saveAs } from "file-saver";

    // 🔥 LOADING STATE (UX penting)
    // 🔥 loading indicator (UX penting saat fetch API)
    const loadingKaryawan = ref(false);

    // 🔥 LIST DATA UNTUK DROPDOWN
    const karyawanOptions = ref([]); // 🔥 state untuk menampung data dropdown karyawan

    // state filter
    const formKaryawanID = ref(null);
    const awal = ref("");
    const akhir = ref("");
    
    // ==============================
    // 🔥 ERROR STATE (INLINE)
    // ==============================
    const errors = ref({
      tanggal: "",
      fetch: "",
      karyawan_id: "",
    });

    const fetchDropdownData = async () => {
    try {
      // loading state
      loadingKaryawan.value = true;
      // 🔥 Panggil API
      // Mengambil daftar semua karyawan
      const resKaryawan = await karyawanService.getAllKaryawan();

      // 🔥 SAFE ACCESS (anti error null / undefined)
      // - resKaryawan?.data?.data → optional chaining (hindari crash jika undefined)
      // - ?? [] → jika hasilnya null / undefined, fallback ke array kosong
      // Kenapa penting?
      // Karena backend kadang return:
      // { data: null } → ini akan bikin .map() error kalau tidak di-handle
      const list = resKaryawan?.data?.data ?? [];

      //console.log("List",list)

      // 🔥 Mapping ke format Vuetify
      // v-autocomplete BUTUH struktur:
      // { title: "yang ditampilkan", value: "yang disimpan" }
      karyawanOptions.value = list.map(k => ({
        title: k.nama,
        value: k.id,
      }));

    } catch (err) {
      console.error(err);
      // 🔥 Fallback UI
      karyawanOptions.value = []; // 🔥 fallback biar UI tidak rusak
      // Set ke array kosong supaya:
      // - tidak crash
      // - autocomplete tetap render normal (meskipun kosong)
    } finally {
      // 🔥 matikan loading
      loadingKaryawan.value = false;
    }
  };

    // watcher tanggal awal dan akhir, auto-clear error saat user input
    watch([awal, akhir], () => {
      errors.value.tanggal = "";
    });

    watch(formKaryawanID, (val) => {
      // kalau autocomplete di-clear
      if (!val) {
        items.value = [];        // 🔥 kosongkan tabel
        errors.value.fetch = ""; // opsional: hilangkan error juga
      }
    });

    // default value (UX bagus)
    onMounted(() => {
        awal.value = dayjs().startOf("month").format("YYYY-MM-DD"); // tglAwal diset ke tgl awal bulan
        akhir.value = dayjs().format("YYYY-MM-DD"); // tglAkhir diset ke tgl hari ini
        // load data karyawan ke autocomplete
        fetchDropdownData();
    }); 

    // state data
    const loading = ref(false);
    const items = ref([]);     // 🔥 untuk tabel (ARRAY)

    // table header
    const headers = [
        { title:"Tanggal", key: "tanggal" },
        { title:"Waktu Masuk", key: "waktu_masuk" },
        { title:"Status", key: "terlambat" },
        { title:"Waktu Pulang", key: "waktu_pulang" },
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
          fetch: "",
          karyawan_id: "",
        };

        // validasi guard (tidak terlalu perlu) karena sudah ada di rule tombol filter
        if (!formKaryawanID.value) {
          errors.value.karyawan_id = "Karyawan wajib dipilih";
          return;
        }

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

            // call 2 api pakai promise
            // const [laporanRes] = await Promise.all([
            //     // api laporan per periode per karyawan per periode
            //     presensiService.getLaporanKaryawanPerPeriode({
            //         id: formKaryawanID.value,
            //         awal : awal.value,
            //         akhir : akhir.value,
            //     }),
            // ]);

            // kalau cuma call 1 api
            // 🔥 SIMPLE: langsung call API
            const laporanRes = await presensiService.getLaporanKaryawanPerPeriode({
              id: formKaryawanID.value,
              awal: awal.value,
              akhir: akhir.value,
            });

            // ambil data dari response json
            items.value = laporanRes.data.data || []; // jika data null, maka kirim map kosong
            // console.log("LAPORAN:", laporanRes.data.data);
            // console.log("TYPE:", typeof laporanRes.data.data);
          } catch (err) {
            console.error(err);
            items.value = [];
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
      // 🔄 reset input
      formKaryawanID.value = null;

      // 🔄 reset tanggal ke default (sama seperti onMounted)
      awal.value = dayjs().startOf("month").format("YYYY-MM-DD");
      akhir.value = dayjs().format("YYYY-MM-DD");

      // 🔄 reset data tabel
      items.value = [];

      // 🔄 reset error
      errors.value = {
        tanggal: "",
        fetch: "",
        karyawan_id: ""
      };
    };

    // 🔥 Ambil object karyawan yang dipilih berdasarkan ID
    // Kenapa pakai computed?
    // - Supaya reactive → otomatis update saat:
    //   1. formKaryawanID berubah
    //   2. karyawanOptions berubah (misalnya setelah fetch API)
    const selectedKaryawan = computed(() => {
      return karyawanOptions.value.find(
        // 🔥 Samakan tipe data pakai String
        // karena ID dari backend kadang number, kadang string
        k => String(k.value) === String(formKaryawanID.value)
      );
    });

    // 🔥 Ambil nama karyawan dari object selectedKaryawan
    // - Kalau belum pilih → tampilkan "-"
    const namaKaryawan = computed(() => {
      return selectedKaryawan.value?.title || "-";
    });

    // ====EXPORT EXCEL====
    const exportExcel = () => {
      // 🔥 guard: kalau belum ada data
      if (!items.value.length) {
        errors.value.fetch = "Tidak ada data untuk diexport";
        return;
      }

      // 🔥 mapping data → format Excel
      const data = items.value.map(item => ({
        "Tanggal": item.tanggal,
        "Waktu Masuk": item.waktu_masuk,
        // "Status": item.terlambat, -> menghasilkan TRUE or FALSE
        // 🔥 mapping boolean → text
        "Status": item.terlambat ? "Terlambat" : "OK", // jika bernilai TRUE -> Terlambat
        "Waktu Pulang": item.waktu_pulang,
      }));

      // 🔥 buat worksheet, mulai dari cell A5 agar tidak tertimpa judul, nama karyawan, dan periode
      const worksheet = XLSX.utils.json_to_sheet(data, { origin: "A5" });

      // Tuliskan judul, nama karyawan, dan periode
      XLSX.utils.sheet_add_aoa(
        worksheet,
        [
          ["Laporan Presensi Karyawan"],
          [`Nama Karyawan: ${namaKaryawan.value}`], // 🔥 ambil value dari computed
          [`Periode: ${awal.value} s/d ${akhir.value}`],
          []
        ],
        { origin: "A1" }
      );

      // auto width column
      worksheet["!cols"] = [
        { wch: 12 }, // Tanggal
        { wch: 12 }, // Waktu Masuk
        { wch: 12 }, // Status
        { wch: 12 }, // Waktu Pulang
      ];

      // 🔥 buat workbook
      const workbook = XLSX.utils.book_new();
      XLSX.utils.book_append_sheet(workbook, worksheet, "Laporan");

      // 🔥 nama file dinamis
      const fileName = `Laporan_Presensi_${namaKaryawan.value}_${awal.value}_sd_${akhir.value}.xlsx`;

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
            Laporan Presensi Per Karyawan
        </div>
      </v-card-title>

      <div class="d-flex align-start gap-3 flex-wrap" style="gap: 12px;">
          <!-- karyawan bisa banyak, jadi pakai autocomplete agar tidak makan space -->
          <v-autocomplete
            v-model="formKaryawanID"
            :items="karyawanOptions"
            item-title="title"
            item-value="value"
            label="Pilih Karyawan"
            placeholder="Cari nama karyawan..."
            density="compact"
            variant="outlined"
            :loading="loadingKaryawan"
            autocomplete="off"
            clearable

            :error="!!errors.karyawan_id"
            :error-messages="errors.karyawan_id"
          />
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
              :disabled="!awal || !akhir || !formKaryawanID"
              @click="fetchLaporan"
              height="40">
              <v-icon start>mdi-filter</v-icon>
              Filter
          </v-btn>

          <!-- 🔥 BUTTON CLEAR -->
          <v-btn
              color="warning"
              :loading="loading"
              :disabled="!awal && !akhir || !formKaryawanID"
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

        <template #item.terlambat="{ item }">
          <v-chip
              :color="item.terlambat ? 'red' : 'green'"
              size="small"
              variant="flat"
          >
              {{ item.terlambat ? 'Terlambat' : 'OK' }}
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