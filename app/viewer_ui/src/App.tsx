import {CssBaseline} from "@mui/material";
import {BrowserRouter} from "react-router-dom";
import DphotoTheme from "./components/dphoto.theme";
import {SecurityIntegration} from "./core/application";
import GeneralRouter from "./pages/general.router";

const App = () => {
  // TODO - add MUI theme
  // TODO - add React error boundary
  return (
    <div className="App">
      <DphotoTheme>
        <CssBaseline/>
        <BrowserRouter>
          <SecurityIntegration>
            <GeneralRouter/>
          </SecurityIntegration>
        </BrowserRouter>
      </DphotoTheme>
    </div>
  )
}

export default App;
