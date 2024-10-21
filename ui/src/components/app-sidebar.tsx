import { Settings, MapPin } from "lucide-react"

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "./ui/collapsible";

const cities = [
    {
        title: "Delhi",
        url: "/city/dehli"
    },
    {
        title: "Mumbai",
        url: "/city/mumbai"
    }
]

export function AppSidebar() {
  return (
    <Sidebar>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Realtime Weather</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
                <Collapsible defaultOpen className="group/collapsible">
                 <SidebarMenuItem>
                 <CollapsibleTrigger asChild>
                 <SidebarMenuButton asChild>
                    <div>
                      <MapPin />
                      <span>Cities</span>
                    </div>
                  </SidebarMenuButton>
                </CollapsibleTrigger>
             <CollapsibleContent>
                    <SidebarMenuSub>
                    {cities.map((city) => (
                <SidebarMenuSubItem key={city.title}>
                  <SidebarMenuSubButton asChild>
                    <a href={city.url}>
                      <span>{city.title}</span>
                    </a>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              ))}
                    </SidebarMenuSub>
            </CollapsibleContent>
                </SidebarMenuItem>
            </Collapsible>
            </SidebarMenu>
            <SidebarMenu>
                <SidebarMenuItem key={0}>
                  <SidebarMenuButton asChild>
                    <a href={"/settings"}>
                      <Settings />
                      <span>Settings</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  )
}
