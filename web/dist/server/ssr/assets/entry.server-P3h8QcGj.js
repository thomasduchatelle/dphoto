import { r as reactExports, S as Slot, C as Children } from "../index.js";
import "../__vite_rsc_assets_manifest.js";
import "node:async_hooks";
const RouterContext = reactExports.createContext(null);
const notAvailableInServer = (name) => () => {
  throw new Error(`${name} is not in the server`);
};
function renderError(message) {
  return reactExports.createElement("html", null, reactExports.createElement("body", null, reactExports.createElement("h1", null, message)));
}
class ErrorBoundary extends reactExports.Component {
  constructor(props) {
    super(props);
    this.state = {};
  }
  static getDerivedStateFromError(error) {
    return {
      error
    };
  }
  render() {
    if ("error" in this.state) {
      if (this.state.error instanceof Error) {
        return renderError(this.state.error.message);
      }
      return renderError(String(this.state.error));
    }
    return this.props.children;
  }
}
const getRouteSlotId = (path) => "route:" + decodeURI(path);
const MOCK_ROUTE_CHANGE_LISTENER = {
  on: () => notAvailableInServer("routeChange:on"),
  off: () => notAvailableInServer("routeChange:off")
};
function INTERNAL_ServerRouter({ route, httpstatus }) {
  const routeElement = reactExports.createElement(Slot, {
    id: getRouteSlotId(route.path)
  });
  const rootElement = reactExports.createElement(Slot, {
    id: "root"
  }, reactExports.createElement("meta", {
    name: "httpstatus",
    content: `${httpstatus}`
  }), routeElement);
  return reactExports.createElement(reactExports.Fragment, null, reactExports.createElement(RouterContext, {
    value: {
      route,
      changeRoute: notAvailableInServer("changeRoute"),
      prefetchRoute: notAvailableInServer("prefetchRoute"),
      routeChangeEvents: MOCK_ROUTE_CHANGE_LISTENER,
      fetchingSlices: /* @__PURE__ */ new Set()
    }
  }, rootElement));
}
const export_847a2b1045ef = {
  Children,
  Slot
};
const export_0f591c01fa0d = {
  ErrorBoundary,
  INTERNAL_ServerRouter
};
export {
  export_0f591c01fa0d,
  export_847a2b1045ef
};
