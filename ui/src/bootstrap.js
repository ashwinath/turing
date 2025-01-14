import React from "react";
import ReactDOM from "react-dom";
import "./assets/style.scss";
import * as serviceWorker from "./serviceWorker";
import App from "./App";
import * as Sentry from "@sentry/browser";
import { ConfigProvider, useConfig } from "./config";

const SentryApp = ({ children }) => {
  const {
    sentryConfig: { dsn, environment, tags },
  } = useConfig();

  Sentry.init({ dsn, environment });
  Sentry.setTags(tags);

  return children;
};

const TuringUI = () => (
  <ConfigProvider>
    <SentryApp>
      <App />
    </SentryApp>
  </ConfigProvider>
);

ReactDOM.render(TuringUI(), document.getElementById("root"));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
