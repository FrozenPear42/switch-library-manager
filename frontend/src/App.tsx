import styles from "./App.module.css";
import Footer from "./components/Footer/Footer";
import Header from "./components/Header/Header";
import LoadingPanel from "./components/LoadingPanel/LoadingPanel";
import Menu from "./components/Menu/Menu";
import Catalog from "./pages/GameCatalog/Catalog";
import Library from "./pages/Library/Library";

export default function App() {
  return (
    <div className={styles.app}>
      <div className={styles.header}>
        <Header></Header>
      </div>
      <div className={styles.menu}>
        <Menu></Menu>
      </div>
      <div className={styles.content}>
        <Catalog></Catalog>
      </div>
      <div className={styles.footer}>
        <Footer></Footer>
      </div>
    </div>
  );

  // return <LoadingPanel></LoadingPanel>;
}
