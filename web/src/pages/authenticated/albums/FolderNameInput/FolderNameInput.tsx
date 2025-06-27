import React from 'react';
import {InputAdornment, Switch, TextField, Tooltip} from "@mui/material";

export interface FolderNameInputProps {
    useCustomFolderName: boolean;
    value: string;
    placeholder: string;
    disabled?: boolean;
    onEnabledChange: (enabled: boolean) => void;
    onValueChange: (value: string) => void;
    tooltip?: string;
    error?: string;
}

export function FolderNameInput({
                                    useCustomFolderName,
    value,
    disabled = false,
    onEnabledChange,
    onValueChange,
    tooltip = "The name of the physical folder name is generated from the date and the name; but can be overridden.",
    error
}: FolderNameInputProps) {
    return (
        <TextField
            variant="outlined"
            label={!useCustomFolderName && !value ? "Automatically generate a folder name" : "Folder name"}
            disabled={!useCustomFolderName || disabled}
            value={value}
            onChange={(event) => onValueChange(event.target.value)}
            error={!!error}
            helperText={error}
            fullWidth
            InputProps={{
                endAdornment: (
                    <InputAdornment position="end">
                        <Tooltip title={tooltip}>
                            <Switch
                                checked={!useCustomFolderName}
                                disabled={disabled}
                                onChange={(event) => onEnabledChange(!event.target.checked)}
                                size="medium"
                            />
                        </Tooltip>
                    </InputAdornment>
                ),
            }}
        />
    );
}
