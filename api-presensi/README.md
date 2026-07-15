## Setup Backend

Masuk ke folder backend

```bash
cd api-presensi
```

Salin environment

```bash
cp .env.example .env.docker
```

Edit

```bash
nano .env.docker
```

Sesuaikan password database dan JWT.

Kembali ke root project

```bash
cd ..
```

Jalankan Docker

```bash
docker compose up -d --build
```