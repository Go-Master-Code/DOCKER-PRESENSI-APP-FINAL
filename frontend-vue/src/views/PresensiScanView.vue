<script setup>
    import {ref, watch, onMounted, onUnmounted, onBeforeUnmount} from "vue";
    import dayjs from "dayjs";
    import * as presensiService from "@/services/presensiService";

    // import logo
    import logo from "@/assets/logo.png";

    // ==============================
    // 🖥️ FULLSCREEN MODE
    // ==============================
    const enterFullscreen = async () => {
        const el = document.documentElement;

        try {
            await el.requestFullscreen();
            isFullscreen.value = true;
        } catch (err) {
            console.warn("Fullscreen gagal:", err);
        }
    };
    // ==============================
    // 🎨 STATUS COLOR (UI LOGIC)
    // ==============================
    const getStatusColor = () => {
        if (!result.value) return "grey";
        if (result.value.error) return "error"; // ❌ error
        if (result.value.status === "pulang") return "primary"; // 🔵 pulang
        if (result.value.status === "masuk") {
            return result.value.terlambat ? "error" : "success"; // ⚠️ / ✅
        }
        return "grey";
    };

    // ==============================
    // 🔥 STATE WAKTU DAN JAM REALTIME
    // ==============================
    const now = ref(dayjs().format("HH:mm:ss"));       // jam realtime
    const today = ref(dayjs().format("YYYY-MM-DD"));   // 🔥 tanggal untuk backend
    const displayDate = ref(dayjs().format("dddd, DD MMMM YYYY")); // 🔥 tanggal untuk UI

    // state fullscreen
    const isFullscreen = ref(false);

    let interval = null;

    onMounted(() => {
        interval = setInterval(() => {
            now.value = dayjs().format("HH:mm:ss"); // UPDATE jam tiap detik

            // 🔥 update tanggal juga (jaga-jaga lewat tengah malam)
            today.value = dayjs().format("YYYY-MM-DD");
            displayDate.value = dayjs().format("dddd, DD MMMM YYYY");
        }, 1000);
    });

    // cleanup (hindari memory leak)
    onUnmounted(() => {
        clearInterval(interval);
    });

    // supaya tidak ada memory leak saat pindah halaman
    onBeforeUnmount(() => {
        clearTimeout(hideNotifTimeout);
    });
    
    // ==============================
    // 🔥 STATE INPUT
    // ==============================
    const karyawanID = ref(""); // input ID karyawan
    const inputRef = ref(null); // untuk autofocus ulang

    // 🔥 Utility: hanya angka (sanitize input)
    // penting karena user bisa paste string
    const onlyNumber = (value) => {
        return value.replace(/\D/g, ""); // hapus semua selain angka
    };

    // 🔥 BLOCK SEMUA KARAKTER SELAIN ANGKA
    const handleKeydown = (e) => {
    // allow ctrl/cmd (copy, paste, dll)
        if (e.ctrlKey || e.metaKey) return;

        // allow control keys
        const allowed = [
            "Backspace",
            "Delete",
            "ArrowLeft",
            "ArrowRight",
            "Tab"
        ];

        if (allowed.includes(e.key)) return;

        // hanya angka
        if (!/[0-9]/.test(e.key)) {
            e.preventDefault(); // 🔥 BLOCK huruf langsung
        }
    };
    
    // ==============================
    // 🔥 STATE UI
    // ==============================
    const loading = ref(false);
    const result = ref(null); // hasil notifikasi (success / error)

    // ==============================
    // 🔊 FEEDBACK SUARA
    // ==============================

    // 🔥 pre-load audio biar tidak delay saat diputar
    const successSound = new Audio("/success.mp3");
    const errorSound = new Audio("/error.mp3");

    // 🔥 function helper
    const playSound = (type) => {
        if (type === "error") {
            errorSound.currentTime = 0;
            errorSound.play();
        } else {
            successSound.currentTime = 0;
            successSound.play();
        }
    };

    // 🔥 timer auto hide notifikasi
    let hideNotifTimeout = null;

    // ==============================
    // 🔥 SUBMIT PRESENSI
    // ==============================
    const submitPresensi = async () => {
        if (!karyawanID.value) return; // guard kosong

        try {
            loading.value = true;

            // 🔥 CALL API
            const res = await presensiService.createPresensi({
                karyawan_id: karyawanID.value,
                // tanggal dan waktu
                tanggal: today.value, // format YYYY-MM-DD
                waktu_masuk: now.value, // format hh:mm:ss
                waktu_pulang: now.value, // format hh:mm:ss
            });

            const data = res.data.data; // ambil response dari data

            // DETEKSI STATUS
            // kalau waktu_pulang = waktu_masuk -> presensi masuk
            const isPulang = data.waktu_pulang && data.waktu_pulang !== data.waktu_masuk

            // 🔥 SIMPAN HASIL (SUCCESS)
            result.value = {
                nama: data.karyawan_nama,
                status: isPulang? "pulang" : "masuk",
                waktu: isPulang? data.waktu_pulang : data.waktu_masuk,
                terlambat: data.terlambat
            }

            playSound("success");

            // 🔥 reset input setelah sukses
            karyawanID.value = "";

        } catch (err) {
            // 🔥 SIMPAN ERROR
            result.value = {
                error: err.response?.data?.error || "Terjadi kesalahan",
            };
            playSound("error");
            // 🔥 reset input setelah sukses
            karyawanID.value = "";
        } finally {
            loading.value = false;

            // 🔥 autofocus ulang (UX penting untuk scanner)
            setTimeout(() => {
                inputRef.value?.focus();
            }, 100);

            // 🔥 auto hide notifikasi
            // 🔥 reset timeout lama
            clearTimeout(hideNotifTimeout);

            // 🔥 buat timeout baru
            hideNotifTimeout = setTimeout(() => {
                result.value = null;
            }, 3000);
        }
    };

    // ==============================
    // 🔥 ENTER HANDLER (INPUT)
    // ==============================
    const handleEnter = (e) => {
        e.preventDefault();
        submitPresensi();
    };

    // ==============================
    // 🔥 GLOBAL ENTER (UNTUK SCANNER)
    // ==============================
    const handleGlobalEnter = (e) => {
        if (e.key !== "Enter") return;
        if (!karyawanID.value) return;
        submitPresensi();
    };

    // ==============================
    // 🔥 AUTO FOCUS LOOP (KIOSK MODE)
    // ==============================
    let focusInterval = null;

    onMounted(() => {
    // fokus pertama
        setTimeout(() => inputRef.value?.focus(), 100);

        // 🔥 jaga fokus tetap di input
        focusInterval = setInterval(() => {
        inputRef.value?.focus();
        }, 2000);

        // 🔥 global listener
        window.addEventListener("keydown", handleGlobalEnter);
    });

    onBeforeUnmount(() => {
        clearInterval(focusInterval);
        window.removeEventListener("keydown", handleGlobalEnter);
    });

    // ==============================
    // 🔥 FALLBACK AUTO SUBMIT
    // (kalau scanner TIDAK kirim Enter)
    // ==============================
    let typingTimeout = null;

    const handleInput = (e) => {
        // 🔥 ambil value dari input langsung
        let value = e.target.value;

        // 🔥 filter hanya angka + max 10 digit
        value = onlyNumber(value).slice(0, 10);

        // 🔥 set kembali ke state
        karyawanID.value = value;

        // 🔥 fallback auto submit (scanner tanpa enter)
        clearTimeout(typingTimeout);

        // tunggu 300ms → kalau tidak ada input lagi → submit
        typingTimeout = setTimeout(() => {
        if (karyawanID.value.length > 9) {
            submitPresensi();
        }
        }, 300);
    };
</script>

<template>
    <div
        v-if="!isFullscreen"
        class="fullscreen-overlay d-flex justify-center align-center"
    >
        <v-btn
            size="x-large"
            color="primary"
            class="text-h5 px-8 py-6"
            @click="enterFullscreen"
        >
            Mulai Presensi
        </v-btn>
    </div>

    <v-container class="fill-height d-flex justify-center align-center">
        <v-card max-width="520" width="100%" class="pa-6 text-center" elevation="6">
            <!-- logo sekolah -->
            <v-img
                :src="logo"
                width="110"
                class="mx-auto mb-1"
            />
            
            <!-- Judul sistem -->
            <div style="font-size: clamp(28px, 1.5vw, 20px);font-weight: 500;">
                Sistem Presensi Karyawan
            </div>

            <!-- ========================= -->
            <!-- 🔥 TANGGAL -->
            <!-- ========================= -->
            <div style="font-size: clamp(20px, 1.5vw, 20px); color: #666;">
                {{ displayDate }}
            </div>

            <!-- ========================= -->
            <!-- 🔥 JAM REALTIME -->
            <!-- ========================= -->
            <div style="font-size: clamp(36px, 5vw, 72px); font-weight: 700;">
                {{ now }}
            </div>

            <!-- ========================= -->
            <!-- 🔥 INPUT ID -->
            <!-- ========================= -->
            <v-text-field
                ref="inputRef"
                v-model="karyawanID"
                @update:modelValue="val => karyawanID = onlyNumber(val).slice(0,10)"

                label="Tap kartu RFID"
                variant="outlined"
                density="comfortable"
                autofocus
                counter
                maxlength="10"
                inputmode="numeric"
                autocomplete="off"

                @keydown="handleKeydown"
                @input="handleInput"

                @keydown.enter.prevent="handleEnter"
            />

            <!-- smooth transition untuk notif error/success presensi -->
            <v-fade-transition>
                <!-- ========================= -->
                <!-- 🔥 NOTIFIKASI -->
                <!-- ========================= -->
                <v-alert
                    v-if="result"
                    class="mt-4 text-subtitle-1"
                    :color="getStatusColor()"
                    variant="tonal"
                    density="comfortable"
                    border="start"
                >
                <!-- ❌ ERROR -->
                <div v-if="result.error">
                    ❌ {{ result.error }}
                </div>

                <!-- ✅ SUCCESS -->
                <div v-else>
                    <!-- Nama -->
                    <div class="font-weight-bold">
                        {{ result.nama }}
                    </div>

                    <!-- Status -->
                    <div class="font-weight-medium">
                        <span v-if="result.status === 'masuk'" class="text-green">
                        ✅ Masuk
                        </span>

                        <span v-else class="text-blue">
                        🔵 Pulang
                        </span>
                    </div>

                    <!-- Waktu -->
                    <div>
                    ⏰ {{ result.waktu }}
                    </div>

                    <!-- 🔥 TERLAMBAT (HANYA UNTUK MASUK) -->
                    <div v-if="result.status === 'masuk'">
                        <span v-if="result.terlambat" class="font-weight-bold text-red">
                        ⚠️ Terlambat
                        </span>
                        <span v-else class="font-weight-bold text-green">
                        ✔️ Tepat waktu
                        </span>
                    </div>
                </div>
                </v-alert>
            </v-fade-transition>
        </v-card>
    </v-container>
</template>

<style>
    .fullscreen-overlay {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: white;
        z-index: 9999;
    }
</style>