package helper

import (
	"api-presensi/internal/dto"
	"api-presensi/internal/model"
)

func ConvertToDTOKaryawanPlural(karyawan []model.Karyawan) []dto.KaryawanResponse {
	// buat var untuk tampung data dto
	var karyawanDTO []dto.KaryawanResponse
	for _, k := range karyawan {
		karyawanDTO = append(karyawanDTO, dto.KaryawanResponse{
			ID:    k.ID,
			Nama:  k.Nama,
			Aktif: k.Aktif,
		})
	}
	return karyawanDTO
}

func ConvertToDTOKaryawanSingle(karyawan model.Karyawan) dto.KaryawanResponse {
	var karyawanDTO dto.KaryawanResponse
	karyawanDTO.ID = karyawan.ID
	karyawanDTO.Nama = karyawan.Nama
	karyawanDTO.Aktif = karyawan.Aktif
	return karyawanDTO
}

func ConvertToDTOJenisIjinSingle(jenisIjin model.JenisIjin) dto.JenisIjinResponse {
	var jenisIjinDTO dto.JenisIjinResponse
	jenisIjinDTO.ID = jenisIjin.ID
	jenisIjinDTO.Nama = jenisIjin.Nama
	jenisIjinDTO.Aktif = jenisIjin.Aktif
	return jenisIjinDTO
}

func ConvertToDTOJenisIjinPlural(jenisIjin []model.JenisIjin) []dto.JenisIjinResponse {
	var jenisIjinDTO []dto.JenisIjinResponse
	for _, j := range jenisIjin {
		jenisIjinDTO = append(jenisIjinDTO, dto.JenisIjinResponse{
			ID:    j.ID,
			Nama:  j.Nama,
			Aktif: j.Aktif,
		})
	}
	return jenisIjinDTO
}

func ConvertToDTOHariLiburPlural(hariLibur []model.HariLibur) []dto.HariLiburResponse {
	var hariLiburDTO []dto.HariLiburResponse
	for _, hl := range hariLibur {
		hariLiburDTO = append(hariLiburDTO, dto.HariLiburResponse{
			ID:         hl.ID,
			Tanggal:    hl.Tanggal.Format("2006-01-02"), // parsing format tanggal yyyy-mm-dd sesuai format mysql
			Hari:       hl.Tanggal.Weekday().String(),   // mendapatkan nama hari dari tanggal
			Keterangan: hl.Keterangan,
		})
	}
	return hariLiburDTO
}

func ConvertToDTOHariLiburSingle(hariLibur model.HariLibur) dto.HariLiburResponse {
	var hariLiburDTO dto.HariLiburResponse
	hariLiburDTO.ID = hariLibur.ID
	hariLiburDTO.Tanggal = hariLibur.Tanggal.Format("2006-01-02")
	hariLiburDTO.Hari = hariLibur.Tanggal.Weekday().String()
	hariLiburDTO.Keterangan = hariLibur.Keterangan
	return hariLiburDTO
}

func ConvertToDTOIjinKaryawanSingle(ijin model.IjinKaryawan) dto.IjinKaryawanResponse {
	var ijinDTO dto.IjinKaryawanResponse
	ijinDTO.ID = ijin.ID
	ijinDTO.Tanggal = ijin.Tanggal.Format("2006-01-02")
	ijinDTO.KaryawanID = ijin.KaryawanID
	ijinDTO.KaryawanNama = ijin.Karyawan.Nama
	ijinDTO.JenisIjinID = ijin.JenisIjinID
	ijinDTO.JenisIjinNama = ijin.JenisIjin.Nama
	ijinDTO.Keterangan = ijin.Keterangan
	return ijinDTO
}

func ConvertToDTOIjinKaryawanPlural(ijin []model.IjinKaryawan) []dto.IjinKaryawanResponse {
	var ijinDTO []dto.IjinKaryawanResponse
	for _, i := range ijin {
		ijinDTO = append(ijinDTO, dto.IjinKaryawanResponse{
			ID:            i.ID,
			Tanggal:       i.Tanggal.Format("2006-01-02"),
			KaryawanID:    i.KaryawanID,
			KaryawanNama:  i.Karyawan.Nama,
			JenisIjinID:   i.JenisIjinID,
			JenisIjinNama: i.JenisIjin.Nama,
			Keterangan:    i.Keterangan,
		})
	}
	return ijinDTO
}

func ConvertToDTOPresensiSingle(presensi model.Presensi) dto.PresensiResponse {
	var presensiDTO dto.PresensiResponse
	presensiDTO.ID = presensi.ID
	presensiDTO.KaryawanID = presensi.KaryawanID
	presensiDTO.KaryawanNama = presensi.Karyawan.Nama
	presensiDTO.Tanggal = presensi.Tanggal.Format("2006-01-02")
	presensiDTO.WaktuMasuk = presensi.WaktuMasuk
	presensiDTO.WaktuPulang = presensi.WaktuPulang
	var terlambat bool
	// status terlambat / tidak
	batasWaktuTerlambat := "07:00:00"
	if presensi.WaktuMasuk > batasWaktuTerlambat {
		terlambat = true
	} else {
		terlambat = false
	}
	presensiDTO.Terlambat = terlambat

	return presensiDTO
}

func ConvertToDTOPresensiPlural(presensi []model.Presensi) []dto.PresensiResponse {
	batasWaktuTerlambat := "07:00:00"
	var presensiDTO []dto.PresensiResponse
	var terlambat bool
	for _, p := range presensi {
		if p.WaktuMasuk > batasWaktuTerlambat {
			terlambat = true
		} else {
			terlambat = false
		}
		presensiDTO = append(presensiDTO, dto.PresensiResponse{
			ID:           p.ID,
			KaryawanID:   p.KaryawanID,
			KaryawanNama: p.Karyawan.Nama,
			Tanggal:      p.Tanggal.Format("2006-01-02"),
			WaktuMasuk:   p.WaktuMasuk,
			WaktuPulang:  p.WaktuPulang,
			Terlambat:    terlambat,
		})
	}
	return presensiDTO
}

func ConvertToDTORolePlural(role []model.Role) []dto.RoleResponse {
	var roleDTO []dto.RoleResponse
	for _, r := range role {
		roleDTO = append(roleDTO, dto.RoleResponse{
			ID:        r.ID,
			Nama:      r.Nama,
			Deskripsi: r.Deskripsi,
		})
	}
	return roleDTO
}

func ConvertToDTORoleSingle(role model.Role) dto.RoleResponse {
	var roleDTO dto.RoleResponse
	roleDTO.ID = role.ID
	roleDTO.Nama = role.Nama
	roleDTO.Deskripsi = role.Deskripsi
	return roleDTO
}

func ConvertToDTOUserPlural(user []model.User) []dto.UserResponse {
	var userDTO []dto.UserResponse
	for _, u := range user {
		userDTO = append(userDTO, dto.UserResponse{
			ID: u.ID,
			//Email:    u.Email,
			Username: u.Username,
			Password: u.Password,
			RoleID:   u.RoleID,
			RoleNama: u.Role.Nama,
		})
	}
	return userDTO
}

func ConvertToDTOUserSingle(user model.User) dto.UserResponse {
	var userDTO dto.UserResponse
	userDTO.ID = user.ID
	//userDTO.Email = user.Email
	userDTO.Username = user.Username
	userDTO.Password = user.Password
	userDTO.RoleID = user.RoleID
	userDTO.RoleNama = user.Role.Nama
	return userDTO
}

func ConvertToDTOLogSingle(log model.Log) dto.LogRequestAndResponse {
	var logDTO dto.LogRequestAndResponse
	logDTO.UserID = log.UserID
	logDTO.Endpoint = log.Endpoint
	logDTO.IPAddress = log.IPAddress
	logDTO.Method = log.Method
	logDTO.UserAgent = log.UserAgent
	return logDTO
}
