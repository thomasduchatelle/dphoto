import {Route, Routes} from "react-router-dom"
import AlbumPage from "./album-content/album-content.page";

const AuthenticatedRouter = () => {
  return (
    <Routes>
      <Route path='/albums' element={<AlbumPage/>}/>
      <Route path='/albums/:owner/:album' element={<AlbumPage/>}/>
      {/*<Route path='*' element={<Navigate to={'/albums'}/>}/>*/}
    </Routes>
  )
}

export default AuthenticatedRouter
