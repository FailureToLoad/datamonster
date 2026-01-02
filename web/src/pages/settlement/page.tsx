import {
  ArchiveIcon,
  PersonIcon,
  type IconProps,
  HourglassMediumIcon,
  TentIcon,
  SignOutIcon,
} from "@phosphor-icons/react";
import { Link, Outlet, useLocation, useNavigate } from "react-router";
import styles from "./page.module.css";

function useNavItems() {
  const { pathname } = useLocation();
  const navigate = useNavigate();
  const timelineKey = "timeline";
  const populationKey = "population";
  const storageKey = "storage";

  const getProps = (active: boolean): IconProps => {
    const props: IconProps = {
      size: 32,
    };
    if (active) {
      props.weight = "fill";
      props.className = "text-primary";
    }
    return props;
  };

  const handleLogout = async () => {
    await fetch("/api/auth/logout", { credentials: "include" });
    navigate("/");
  };

  return {
    timelineKey,
    populationKey,
    storageKey,
    pathname,
    getProps,
    handleLogout,
  };
}

function getPageTitle(pathname: string): string {
  if (pathname.includes("timeline")) return "Timeline";
  if (pathname.includes("population")) return "Population";
  if (pathname.includes("storage")) return "Storage";
  return "Settlement";
}

function getActiveIcon(pathname: string): React.ReactNode {
  if (pathname.includes("timeline")) return <HourglassMediumIcon size={24} weight="fill" />;
  if (pathname.includes("population")) return <PersonIcon size={24} weight="fill" />;
  if (pathname.includes("storage")) return <ArchiveIcon size={24} weight="fill" />;
  return <HourglassMediumIcon size={24} />;
}

function NavBar() {
  const { timelineKey, populationKey, storageKey, pathname, handleLogout } = useNavItems();
  const pageTitle = getPageTitle(pathname);

  return (
    <div className={styles.navWrapper}>
      <div className={styles.navbar}>
        <div className={styles.navbarStart}>
          <div className={styles.tooltip} data-tip="Settlements">
            <Link to="/settlements" className={styles.navLink} aria-label="Settlements">
              <TentIcon size={24} />
            </Link>
          </div>
        </div>
        <div className={styles.navbarCenter}>
          <div className={styles.dropdown}>
            <div tabIndex={0} role="button" className={styles.navButton}>
              {getActiveIcon(pathname)}
              {pageTitle}
            </div>
            <ul tabIndex={-1} className={styles.dropdownMenu}>
              {!pathname.includes(timelineKey) && (
                <li>
                  <Link to={timelineKey}>
                    <HourglassMediumIcon size={24} />
                    Timeline
                  </Link>
                </li>
              )}
              {!pathname.includes(populationKey) && (
                <li>
                  <Link to={populationKey}>
                    <PersonIcon size={24} />
                    Population
                  </Link>
                </li>
              )}
              {!pathname.includes(storageKey) && (
                <li>
                  <Link to={storageKey}>
                    <ArchiveIcon size={24} />
                    Storage
                  </Link>
                </li>
              )}
            </ul>
          </div>
        </div>
        <div className={styles.navbarEnd}>
          <div className={styles.tooltip} data-tip="Logout">
            <button onClick={handleLogout} className={styles.navLink} aria-label="Logout">
              <SignOutIcon size={24} />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export function SettlementPage() {
  return (
    <div className={styles.page}>
      <NavBar />
      <div className={styles.content}>
        <Outlet />
      </div>
    </div>
  );
}
