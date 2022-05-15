import {CssBaseline} from "@mui/material";
import {BrowserRouter} from "react-router-dom";
import DphotoTheme from "./components/DPhotoTheme";
import {AppConfigIntegration, SecurityIntegration} from "./core/application";
import GeneralRouter from "./pages/GeneralRouter";

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
