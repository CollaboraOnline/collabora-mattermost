import { FileInfo } from '@mattermost/types/files';

export interface PluginRegistry {
  registerPostTypeComponent(typeName: string, component: React.ElementType);
  registerFilePreviewComponent(condition: (fileInfo: FileInfo) => boolean, component: React.ElementType);
  registerFileDropdownMenuAction(
    shouldShowPreview: (store: Store<GlobalState, AnyAction>, fileInfo: FileInfo) => boolean,
    name: string,
    action: (fileInfo: FileInfo) => void
  );
  registerFileUploadMethod(component: React.Element, action: () => void, name: string);

  registerChannelHeaderButtonAction(icon: React.Element, action: () => void, dropdownText: string, tooltipText: string);
  registerChannelIntroButtonAction(icon: React.Element, action: () => void, tooltipText: string);
  registerCustomRoute(route: string, component: React.ElementType);
  registerProductRoute(route: string, component: React.ElementType);
  unregisterComponent(componentId: string);
  registerProduct(
    baseURL: string,
    switcherIcon: string,
    switcherText: string,
    switcherLinkURL: string,
    mainComponent: React.ElementType,
    headerCentreComponent: React.ElementType,
    headerRightComponent?: React.ElementType,
    showTeamSidebar: boolean
  );
  registerPostWillRenderEmbedComponent(match: (embed: { type: string; data: any }) => void, component: any, toggleable: boolean);
  registerWebSocketEventHandler(event: string, handler: (e: any) => void);
  unregisterWebSocketEventHandler(event: string);
  registerAppBarComponent(iconURL: string, action: (channel: Channel, member: ChannelMembership) => void, tooltipText: React.ReactNode);
  registerRightHandSidebarComponent(component: React.ElementType, title: React.Element);
  registerRootComponent(component: React.ElementType);
  registerInsightsHandler(handler: (timeRange: string, page: number, perPage: number, teamId: string, insightType: string) => void);
  registerSiteStatisticsHandler(handler: () => void);
  registerActionAfterChannelCreation(component: React.Element);
  registerSlashCommandWillBePostedHook(
    slashCommandWillBePostedHook: (message: any, contextArgs: any) => Promise<{ message: any; args: any }> | { message: any; args: any }
  );
  registerReducer(reducer);

  // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference
}
