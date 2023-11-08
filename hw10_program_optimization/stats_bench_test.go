package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	b.Helper()
	b.StopTimer()

	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()

	data, _ := r.File[0].Open()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "biz")
	}
	b.StopTimer()
}
