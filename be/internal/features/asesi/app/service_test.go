package app

import (
	"context"
	"testing"
	"time"

	"lsp/internal/features/asesi/domain"
)

// fakeStorage mencatat object key yang dihapus lewat Delete.
type fakeStorage struct{ deleted []string }

func (f *fakeStorage) PresignedPutURL(context.Context, string, string, time.Duration) (string, error) {
	return "", nil
}
func (f *fakeStorage) PresignedGetURL(context.Context, string, time.Duration) (string, error) {
	return "", nil
}
func (f *fakeStorage) PublicURL(string) string { return "" }
func (f *fakeStorage) Delete(_ context.Context, objectKey string) error {
	f.deleted = append(f.deleted, objectKey)
	return nil
}
func (f *fakeStorage) Exists(context.Context, string) (bool, error) { return false, nil }

// fakeRepo mengembalikan profil lama/baru yang sudah ditentukan.
type fakeRepo struct {
	previous domain.Profile
	updated  domain.Profile
}

func (f *fakeRepo) Create(context.Context, domain.Profile) error { return nil }
func (f *fakeRepo) Patch(context.Context, domain.Patch) (domain.Profile, error) {
	return f.updated, nil
}
func (f *fakeRepo) FindByID(context.Context, string) (domain.Profile, error) { return f.previous, nil }
func (f *fakeRepo) FindByUserID(context.Context, string) (domain.Profile, error) {
	return domain.Profile{}, nil
}
func (f *fakeRepo) List(context.Context, domain.ListParams) (domain.ListResult, error) {
	return domain.ListResult{}, nil
}
func (f *fakeRepo) Delete(context.Context, string) error                      { return nil }
func (f *fakeRepo) CreateOccupation(context.Context, domain.Occupation) error { return nil }
func (f *fakeRepo) PatchOccupation(context.Context, domain.OccupationPatch) (domain.Occupation, error) {
	return domain.Occupation{}, nil
}
func (f *fakeRepo) FindOccupationByID(context.Context, string) (domain.Occupation, error) {
	return domain.Occupation{}, nil
}
func (f *fakeRepo) FindOccupationByUserID(context.Context, string) (domain.Occupation, error) {
	return domain.Occupation{}, nil
}
func (f *fakeRepo) UpsertOccupationByUser(context.Context, string, domain.Occupation) (domain.Occupation, error) {
	return domain.Occupation{}, nil
}
func (f *fakeRepo) ListOccupations(context.Context, domain.ListParams) (domain.OccupationListResult, error) {
	return domain.OccupationListResult{}, nil
}
func (f *fakeRepo) DeleteOccupation(context.Context, string) error { return nil }

func strp(s string) *string { return &s }

func TestRemoveReplacedFiles(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		previous domain.Profile
		updated  domain.Profile
		want     []string
	}{
		{
			name:     "berkas diganti -> lama dihapus",
			previous: domain.Profile{Photo: strp("old/photo.jpg")},
			updated:  domain.Profile{Photo: strp("new/photo.jpg")},
			want:     []string{"old/photo.jpg"},
		},
		{
			name:     "nilai tidak berubah -> tidak dihapus",
			previous: domain.Profile{Photo: strp("same.jpg")},
			updated:  domain.Profile{Photo: strp("same.jpg")},
			want:     nil,
		},
		{
			name:     "berkas lama kosong -> tidak dihapus",
			previous: domain.Profile{Photo: nil},
			updated:  domain.Profile{Photo: strp("new.jpg")},
			want:     nil,
		},
		{
			name: "beberapa berkas diganti sekaligus",
			previous: domain.Profile{
				Photo: strp("a.jpg"), KTPFile: strp("b.pdf"), CVFile: strp("c.pdf"), NPWPFile: strp("d.pdf"),
			},
			updated: domain.Profile{
				Photo: strp("a2.jpg"), KTPFile: strp("b.pdf"), CVFile: strp("c2.pdf"), NPWPFile: strp("d2.pdf"),
			},
			want: []string{"a.jpg", "c.pdf", "d.pdf"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeStorage{}
			svc := &Service{files: store}
			svc.removeReplacedFiles(ctx, tc.previous, tc.updated)
			if len(store.deleted) != len(tc.want) {
				t.Fatalf("deleted = %v, want %v", store.deleted, tc.want)
			}
			for i := range tc.want {
				if store.deleted[i] != tc.want[i] {
					t.Fatalf("deleted = %v, want %v", store.deleted, tc.want)
				}
			}
		})
	}
}

// Patch harus memanggil penghapusan berkas lama setelah menyimpan nilai baru.
func TestPatchDeletesReplacedFile(t *testing.T) {
	store := &fakeStorage{}
	repo := &fakeRepo{
		previous: domain.Profile{ID: "p1", Photo: strp("uploads/photo/old.jpg")},
		updated:  domain.Profile{ID: "p1", Photo: strp("uploads/photo/new.jpg")},
	}
	svc := NewService(repo, nil, store)

	if _, err := svc.Patch(context.Background(), PatchCommand{ID: "p1", Photo: strp("uploads/photo/new.jpg")}); err != nil {
		t.Fatalf("Patch error: %v", err)
	}
	if len(store.deleted) != 1 || store.deleted[0] != "uploads/photo/old.jpg" {
		t.Fatalf("deleted = %v, want [uploads/photo/old.jpg]", store.deleted)
	}
}

// Tanpa object storage (belum dikonfigurasi), tidak boleh panic dan tidak menghapus apa pun.
func TestRemoveReplacedFilesWithoutStorage(t *testing.T) {
	svc := &Service{}
	svc.removeReplacedFiles(context.Background(),
		domain.Profile{Photo: strp("x.jpg")},
		domain.Profile{Photo: strp("y.jpg")},
	)
}
