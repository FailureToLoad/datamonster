import Archive from "lucide-solid/icons/archive";
import PersonStanding from "lucide-solid/icons/person-standing";
import Hourglass from "lucide-solid/icons/hourglass";
import Tent from "lucide-solid/icons/tent";
import LogOut from "lucide-solid/icons/log-out";
import { A, useLocation } from '@solidjs/router';
import styles from "./page.module.css";
import {type ParentProps } from 'solid-js';

function useNavItems() {
  const { pathname } = useLocation();
  const timelineKey = "timeline";
  const populationKey = "population";
  const storageKey = "storage";

  const handleLogout = async () => {
    await fetch("/auth/logout", { credentials: "include" });
    window.location.href = '/auth/login';
  };

  return {
    timelineKey,
    populationKey,
    storageKey,
    pathname,
    handleLogout,
  };
}

function getPageTitle(pathname: string): string {
  if (pathname.includes("timeline")) return "Timeline";
  if (pathname.includes("population")) return "Population";
  if (pathname.includes("storage")) return "Storage";
  return "Timeline";
}

function getActiveIcon(pathname: string) {
  if (pathname.includes("timeline")) return <Hourglass />;
  if (pathname.includes("population")) return <PersonStanding />;
  if (pathname.includes("storage")) return <LogOut />;
  return <Hourglass size={24} />;
}

function NavBar() {
  const { timelineKey, populationKey, storageKey, pathname, handleLogout } = useNavItems();
  const pageTitle = getPageTitle(pathname);

  return (
    <div class={styles.navWrapper}>
      <div class={styles.navbar}>
        <div class={styles.navbarStart}>
          <div class={styles.tooltip} data-tip="Settlements">
            <A href="/settlements" class={styles.navLink} aria-label="Settlements">
              <Tent size={24} />
            </A>
          </div>
        </div>
        <div class={styles.navbarCenter}>
          <div class={styles.dropdown}>
            <div tabIndex={0} role="button" class={styles.navButton}>
              {getActiveIcon(pathname)}
              {pageTitle}
            </div>
            <ul tabIndex={-1} class={styles.dropdownMenu}>
              {!pathname.includes(timelineKey) && (
                <li>
                  <A href={timelineKey}>
                    <Hourglass size={24} />
                    Timeline
                  </A>
                </li>
              )}
              {!pathname.includes(populationKey) && (
                <li>
                  <A href={populationKey}>
                    <PersonStanding size={24} />
                    Population
                  </A>
                </li>
              )}
              {!pathname.includes(storageKey) && (
                <li>
                  <A href={storageKey}>
                    <Archive size={24} />
                    Storage
                  </A>
                </li>
              )}
            </ul>
          </div>
        </div>
        <div class={styles.navbarEnd}>
          <div class={styles.tooltip} data-tip="Logout">
            <button onClick={handleLogout} class={styles.navLink} aria-label="Logout">
              <LogOut size={24} />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export function SettlementPage(props: ParentProps) {
  return (
    <div class={styles.page}>
      <NavBar />
      <div class={styles.content}>
        {props.children}
      </div>
    </div>
  );
}
