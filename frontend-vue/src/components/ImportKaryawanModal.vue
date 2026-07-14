<script setup>
    import { ref } from "vue";
    import * as karyawanService from "@/services/karyawanService";

    const props = defineProps({
    modelValue: Boolean
    });

    const emit = defineEmits(["update:modelValue", "success"]);

    // state
    const file = ref(null);
    const loading = ref(false);

    // 🔥 feedback UI
    const errorMsg = ref("");
    const result = ref(null);

    // 🔥 close modal
    const closeModal = () => {
    emit("update:modelValue", false);

    // reset state biar fresh saat buka lagi
    file.value = null;
    errorMsg.value = "";
    result.value = null;
    };

    // 🔥 upload handler
    const upload = async () => {
    if (!file.value) {
        errorMsg.value = "Silakan pilih file terlebih dahulu.";
        return;
    }

    const formData = new FormData();
    formData.append("file", file.value);

    try {
        loading.value = true;
        errorMsg.value = "";
        result.value = null;

        const res = await karyawanService.importExcel(formData);

        // 🔥 simpan hasil dari backend
        result.value = res.data;

        emit("success"); //🔥 trigger refresh parent
    } catch (err) {
            console.error(err);
            errorMsg.value =
            err.response?.data?.message ||
            "Gagal import data";
    } finally {
        loading.value = false;
    }
    };
</script>

<template>
  <v-dialog v-model="props.modelValue" max-width="400">
    <v-card>
        <!-- TITLE -->
        <v-card-title class="text-h6">
        Import Data Karyawan
        </v-card-title>

        <!-- CONTENT -->
        <v-card-text>
            <v-file-input
            v-model="file"
            label="Upload File Excel"
            accept=".xlsx, .xls"
            variant="outlined"
            density="compact"
            prepend-icon="mdi-file-excel"
            show-size
            clearable
        />

            <!-- ERROR -->
            <v-alert
            v-if="errorMsg"
            type="error"
            density="compact"
            class="mt-3"
            >
            {{ errorMsg }}
            </v-alert>

            <!-- SUCCESS RESULT -->
            <v-alert
                v-if="result && result.inserted_data > 0"
                type="success"
                density="compact"
                class="mt-3"
            >
                Berhasil import {{ result.inserted_data || 0 }} data
            </v-alert>

            <!-- FAILED  RESULT -->
            <v-alert
                v-if="result && result.error_count > 0"
                type="warning"
                density="compact"
                class="mt-2"
            >
                Gagal import {{ result.error_count || 0 }} data
            </v-alert>

            <!-- LIST ERROR PER ROW -->
            <v-expansion-panels
                v-if="result?.errors?.length"
                class="mt-2"
            >
            <v-expansion-panel>
                <v-expansion-panel-title>
                    Lihat detail gagal import
                </v-expansion-panel-title>

                <v-expansion-panel-text>
                <v-list density="compact">

                    <div v-for="(err, i) in result.errors" :key="i" class="error-row">
  ❌ Baris {{ err.row }} - {{ err.message }}
</div>

                </v-list>
                </v-expansion-panel-text>
            </v-expansion-panel>
            </v-expansion-panels>
        </v-card-text>

        <!-- ACTION -->
        <v-card-actions>
            <v-spacer />
                <v-btn
                    color="red"
                    variant="tonal"
                    size="small"
                    @click="closeModal"
                    >
                    Batal
                </v-btn>
                <v-btn
                    color="primary"
                    variant="elevated"
                    size="small"
                    :loading="loading"
                    :disabled="!file"
                    @click="upload"
                    >
                    Upload
                </v-btn>
        </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
    .error-row {
    font-size: 15px;
    color: #d32f2f; /* text-error */
    margin-bottom: 5px;
    }
</style>