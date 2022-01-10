import {Box, Container, CssBaseline} from "@mui/material";

export default () => {
  return (
    <Container component="div"
               sx={{
                 marginTop: 3,
                 width: '650px',
                 textAlign: "center",
                 margin: '8 auto 0 auto'
               }}>
      <CssBaseline/>
      <img src="/dphoto-fulllogo-medium.png" alt="dphoto-logo"/>
      <Box sx={{margin: 'auto', mt: 12}}>
        <img src='/loading.svg'/>
      </Box>
    </Container>
  )
}