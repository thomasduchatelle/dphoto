import React, {useState} from 'react';
import {EditDatesDialogSelection} from '../../core/catalog/album-edit-dates/selector-editDatesDialog';

export interface EditDatesDialogProps {
    selection: EditDatesDialogSelection;
    onSave: () => Promise<void>;
    onClose: () => void;
}

export function EditDatesDialog({selection, onSave, onClose}: EditDatesDialogProps) {
    const [startDate, setStartDate] = useState(selection.currentStartDate);
    const [endDate, setEndDate] = useState(selection.currentEndDate);

    if (!selection.isOpen) {
        return null;
    }

    const handleSave = async () => {
        await onSave();
    };

    return (
        <div className="edit-dates-dialog">
            <div className="dialog-content">
                <h2>Edit Dates - {selection.albumName}</h2>
                
                <div className="date-inputs">
                    <label>
                        Start Date:
                        <input
                            type="date"
                            value={startDate.toISOString().split('T')[0]}
                            onChange={(e) => setStartDate(new Date(e.target.value))}
                            disabled={selection.isLoading}
                        />
                    </label>
                    
                    <label>
                        End Date:
                        <input
                            type="date"
                            value={endDate.toISOString().split('T')[0]}
                            onChange={(e) => setEndDate(new Date(e.target.value))}
                            disabled={selection.isLoading}
                        />
                    </label>
                </div>

                <div className="dialog-actions">
                    <button 
                        onClick={handleSave} 
                        disabled={selection.isLoading}
                        className="save-button"
                    >
                        {selection.isLoading ? 'Saving...' : 'Save'}
                    </button>
                    <button 
                        onClick={onClose} 
                        disabled={selection.isLoading}
                        className="cancel-button"
                    >
                        Cancel
                    </button>
                </div>
            </div>
        </div>
    );
}
