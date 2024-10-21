import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import {
  createBrowserRouter,
  RouterProvider,
  Link,
} from "react-router-dom";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar"
import App from './App';
import Settings from './Settings';

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />
  },
  {
    path: "settings",
    element: <Settings />
  },
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <SidebarProvider>
      <AppSidebar />
      <main>
        <SidebarTrigger />
        <RouterProvider router={router} />
      </main>
    </SidebarProvider>
  </StrictMode>,
)


