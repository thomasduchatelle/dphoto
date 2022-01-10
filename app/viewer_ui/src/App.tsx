import {BrowserRouter} from "react-router-dom";
import {bootstrap} from "./domain/bootstrap";
import AnonymousRouter from "./pages/anonymous.router";
import LoadingPage from "./pages/layout/loading.page";
import SecurityIntegation from "./pages/layout/security-integation";

// resolve all dependency injection and configuration
bootstrap()

const App = () => {
  return (
    <div className="App">
      <BrowserRouter>
        <SecurityIntegation loading={<LoadingPage/>}>
          <AnonymousRouter/>
        </SecurityIntegation>
      </BrowserRouter>
      {/*<BrowserRouter>*/}
      {/*  <SecuredContent>*/}
      {/*    <Home/>*/}
      {/*  </SecuredContent>*/}
      {/*</BrowserRouter>*/}
    </div>
  )
}

export default App;
