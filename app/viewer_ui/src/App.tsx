import {CssBaseline} from "@mui/material";
import {BrowserRouter} from "react-router-dom";
import DphotoTheme from "./components/dphoto.theme";
import {SecurityIntegration} from "./core/application";
import AppConfigIntegration from "./core/application/app-config.integration";
import GeneralRouter from "./pages/general.router";

const App = () => {
  // TODO - add React error boundary
  return (
    <div className="App">
      <DphotoTheme>
        <CssBaseline/>
        <BrowserRouter>
          <AppConfigIntegration>
            <SecurityIntegration>
              <GeneralRouter/>
            </SecurityIntegration>
          </AppConfigIntegration>
        </BrowserRouter>
      </DphotoTheme>
    </div>
  )
}

export default App;
