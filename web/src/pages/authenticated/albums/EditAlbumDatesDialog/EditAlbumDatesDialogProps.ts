export interface EditAlbumDatesDialogProps {
    isOpen: boolean;
    albumName: string;
    startDate: Date;
    endDate: Date;
    isStartDateAtStartOfDay: boolean;
    isEndDateAtEndOfDay: boolean;
    onClose: () => void;
}
