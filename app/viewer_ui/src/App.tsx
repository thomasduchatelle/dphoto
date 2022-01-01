import {BrowserRouter} from "react-router-dom";
import Home from "./Home";
import SecuredContent from "./security/secured-content.case";

const App = () => {
  return (
    <div className="App">
      <BrowserRouter>
        <SecuredContent>
          <Home/>
        </SecuredContent>
      </BrowserRouter>
    </div>
  )
}

export default App;
