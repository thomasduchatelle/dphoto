import LocalFireDepartmentIcon from "@mui/icons-material/LocalFireDepartment";
import {Box} from "@mui/material";

export const HotIndicator = ({relativeTemperature, count}: {
  relativeTemperature: number
  count: number
}) => {
  return (
    <Box sx={{
      position: 'relative',
      display: 'inline-flex',
      width: '35px',
      height: '35px'
    }}>
      <LocalFireDepartmentIcon fontSize='large' sx={{
        position: 'absolute',
        color: 'rgba(0, 0, 0, 0.04)',
      }}/>
      <Box sx={{
        position: 'absolute',
        bottom: 0,
        height: `${10 + relativeTemperature * 90}%`,
        width: '100%',
        overflow: 'hidden'
      }}>
        <LocalFireDepartmentIcon fontSize='large' sx={{
          color: (theme) => theme.palette.error.main,
          position: 'absolute',
          bottom: 0,
          width: '35px',
          height: '35px',
        }}/>
      </Box>
    </Box>)
}