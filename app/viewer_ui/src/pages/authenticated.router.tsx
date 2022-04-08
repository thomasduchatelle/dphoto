import {Route, Routes} from "react-router-dom"
import Home from "./home";

const AuthenticatedRouter = () => {
  return (
    <Routes>
      <Route path='/*' element={<Home/>}/>
    </Routes>
  )
}

export default AuthenticatedRouter
