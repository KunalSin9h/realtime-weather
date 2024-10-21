import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar"
import App from './App';
import Settings from './Settings';

import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast';
import City from './City';

const queryClient = new QueryClient()


const router = createBrowserRouter([
  {
    path: "/",
    element: <App />
  },
  {
    path: "settings",
    element: <Settings />
  },
  {
    path: "city/:city_name/:city_id",
    element: <City />
  }
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
    <SidebarProvider>
      <AppSidebar />
      <Toaster />
      <main className='h-screen w-screen'>
        <SidebarTrigger />
        <RouterProvider router={router} />
      </main>
    </SidebarProvider>
    </QueryClientProvider>
  </StrictMode>,
)


