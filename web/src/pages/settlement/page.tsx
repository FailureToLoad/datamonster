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


function LeftNav() {
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

  return (
    <div className={styles.navContainer}>
      <ul className={`menu rounded-box bg-base-300 ${styles.navMenu}`}>
        <div>
          <li>
            <Link to={timelineKey} color="foreground" className="tooltip tooltip-right" data-tip="Timeline">
              <HourglassMediumIcon {...getProps(pathname.includes(timelineKey))} />
            </Link>
          </li>
          <li>
            <Link to={populationKey} color="foreground" className="tooltip tooltip-right" data-tip="Population">
              <PersonIcon {...getProps(pathname.includes(populationKey))} />
            </Link>
          </li>
          <li>
            <Link to={storageKey} color="foreground" className="tooltip tooltip-right" data-tip="Storage">
              <ArchiveIcon {...getProps(pathname.includes(storageKey))} />
            </Link>
          </li>
        </div>
        <div>
          <li>
            <Link to="/settlements" className="tooltip tooltip-right" data-tip="Settlements">
              <TentIcon size={32} />
            </Link>
          </li>
          <li>
            <button onClick={handleLogout} className="tooltip tooltip-right" data-tip="Logout">
              <SignOutIcon size={32} />
            </button>
          </li>
        </div>
      </ul>
    </div>
  );
}

export function SettlementPage() {
  return (
    <div className={styles.page}>
      <LeftNav />
      <div className={styles.content}>
        <Outlet />
      </div>
    </div>
  );
}
