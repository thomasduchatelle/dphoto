import {Backdrop, CircularProgress} from "@mui/material";

export default ({loading}: {
  loading: boolean
}) => (
  <Backdrop
    sx={{color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1}}
    open={loading}
  >
    <CircularProgress color="primary" size={100} thickness={2}/>
  </Backdrop>
)
