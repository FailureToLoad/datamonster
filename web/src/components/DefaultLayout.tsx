import {Outlet} from 'react-router';
import styles from './ProtectedLayout.module.css';
import {GlossaryProvider} from '~/context/glossary';

export default function Layout() {
    return (
        <GlossaryProvider>
            <div className={styles.layout}>
                <Outlet />
            </div>
        </GlossaryProvider>
    );
}
