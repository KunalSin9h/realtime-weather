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
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";

type Cities = {
    ID: number;
    Name: string;
}

export function AppSidebar() {
    const [cities, setCities] = useState<Cities[]>([])

  const { data, isLoading } = useQuery({
    queryKey: ['cities'],
    queryFn: async () => {
      const resp = await fetch("http://localhost:5000/cities");
      if (!resp.ok) {
        return null;
      }

      return resp.json();
    }
  })

  useEffect(() => {
    if (data) setCities(data);
  }, [isLoading])

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
                    {cities.length > 0 && cities.map((city, index) => (
                <SidebarMenuSubItem key={index}>
                  <SidebarMenuSubButton asChild>
                    <a href={`city&city_id=${city.ID}`}>
                      <span>{city.Name}</span>
                    </a>
                    
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              ))}

                    <SidebarMenuSubItem>
                    {isLoading ? <p className="text-xs text-gray-400 pl-4">Loading...</p> : cities.length === 0 && <p className="text-xs text-gray-400 pl-4">No Cities</p>}
                    </SidebarMenuSubItem>
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
