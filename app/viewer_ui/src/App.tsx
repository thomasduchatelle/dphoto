import {CssBaseline} from "@mui/material";
import {BrowserRouter} from "react-router-dom";
import {SecurityIntegration} from "./core/application";
import GeneralRouter from "./pages/general.router";

const App = () => {
  // TODO - add MUI theme
  // TODO - add React error boundary
  return (
    <div className="App">
      <CssBaseline/>
      <BrowserRouter>
        <SecurityIntegration>
          <GeneralRouter/>
        </SecurityIntegration>
      </BrowserRouter>
    </div>
  )
}

export default App;
