<script setup>
import { useRouter } from "vue-router";

import {
    useSessionExpired
} from "@/composables/useSessionExpired";

const router = useRouter();

const {
    showSessionExpired,
    closeSessionExpired,
} = useSessionExpired();

// =========================
// 🔥 LOGIN ULANG
// =========================
const handleRelogin = () => {

    // hapus token
    localStorage.removeItem("token");
    localStorage.removeItem("role");

    // tutup dialog
    closeSessionExpired();

    // redirect login
    router.push("/login");
};
</script>

<template>
    <!-- atribut persistent artinya modal tidak bisa di lost focus, tidak bisa ditutup, tidak bisa di esc -->
    <v-dialog
        v-model="showSessionExpired"
        persistent
        width="420"
    >
        <v-card
            rounded="xl"
            class="pa-4"
        >
            <!-- ICON -->
            <div class="text-center mb-4">
                <v-icon
                    size="64"
                    color="warning"
                >
                    mdi-clock-alert-outline
                </v-icon>
            </div>

            <!-- TITLE -->
            <div
                class="text-h5 font-weight-bold text-center mb-2"
            >
                Session Berakhir
            </div>

            <!-- MESSAGE -->
            <div
                class="text-body-1 text-medium-emphasis text-center mb-6"
            >
                Session login Anda telah berakhir.
                Silakan login kembali untuk melanjutkan.
            </div>

            <!-- ACTION -->
            <v-btn
                block
                color="primary"
                size="large"
                @click="handleRelogin"
            >
                Oke
            </v-btn>
        </v-card>
    </v-dialog>
</template>