package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		b.Logf(err.Error())
		return
	}
	defer r.Close()

	data, err := r.File[0].Open()
	if err != nil {
		b.Logf(err.Error())
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "biz")
	}
}
