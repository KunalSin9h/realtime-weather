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

import {
  useQuery,
  useMutation,
  useQueryClient,
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'

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
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
    <SidebarProvider>
      <AppSidebar />
      <main className='h-screen w-screen'>
        <SidebarTrigger />
        <RouterProvider router={router} />
      </main>
    </SidebarProvider>
    </QueryClientProvider>
  </StrictMode>,
)


