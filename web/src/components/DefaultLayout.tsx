import {Outlet} from 'react-router';
import styles from './ProtectedLayout.module.css';

export default function Layout() {
    return (
        <div className={styles.layout}>
            <Outlet />
        </div>
    );
}
