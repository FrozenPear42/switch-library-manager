import { Redirect, Route, Switch } from "wouter";
import styles from "./App.module.css";
import Footer from "./components/Footer/Footer";
import Header from "./components/Header/Header";
import LoadingPanel from "./components/LoadingPanel/LoadingPanel";
import Menu from "./components/Menu/Menu";
import Catalog from "./pages/GameCatalog/Catalog";
import Library from "./pages/Library/Library";
import Files from "./pages/Files/Files";
import { useStartup } from "./hooks/useStartup";
import NUTPage from "./pages/NUTPage/NUTPage";

export default function App() {
  const { state } = useStartup();

  if (state?.running) {
    return <LoadingPanel state={state}></LoadingPanel>;
  }

  return (
    <div className={styles.app}>
      <div className={styles.header}>
        <Header></Header>
      </div>
      <div className={styles.menu}>
        <Menu></Menu>
      </div>
      <div className={styles.content}>
        <Switch>
          <Route path="/catalog">
            <Catalog />
          </Route>
          <Route path="/library">
            <Library />
          </Route>
          <Route path="/files">
            <Files />
          </Route>
          <Route path="/nut">
            <NUTPage />
          </Route>
          <Route path="/settings">Settings</Route>
          <Route>
            <Redirect to="/library"></Redirect>
          </Route>
        </Switch>
      </div>
      <div className={styles.footer}>
        <Footer></Footer>
      </div>
    </div>
  );
}
