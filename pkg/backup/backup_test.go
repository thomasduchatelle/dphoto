package backup

//func TestBackupAcceptance(t *testing.T) {
//	type args struct {
//		owner        ownermodel.Owner
//		volume       SourceVolume
//		optionsSlice []Options
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    CompletionReport
//		wantErr assert.ErrorAssertionFunc
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := Backup(tt.args.owner, tt.args.volume, tt.args.optionsSlice...)
//			if !tt.wantErr(t, err, fmt.Sprintf("Backup(%v, %v, %v)", tt.args.owner, tt.args.volume, tt.args.optionsSlice...)) {
//				return
//			}
//			assert.Equalf(t, tt.want, got, "Backup(%v, %v, %v)", tt.args.owner, tt.args.volume, tt.args.optionsSlice...)
//		})
//	}
//}
