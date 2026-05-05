#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>

extern void madoTrayToggle(void);
extern void madoTrayQuit(void);

@interface MadoTrayMenuTarget : NSObject
- (void)toggle:(id)sender;
- (void)quit:(id)sender;
@end

@implementation MadoTrayMenuTarget
- (void)toggle:(id)sender {
  madoTrayToggle();
}

- (void)quit:(id)sender {
  madoTrayQuit();
}
@end

static NSStatusItem *madoTrayStatusItem;
static MadoTrayMenuTarget *madoTrayTarget;

void MadoTrayCreate(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (madoTrayStatusItem != nil) {
      return;
    }

    madoTrayTarget = [[MadoTrayMenuTarget alloc] init];
    madoTrayStatusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    madoTrayStatusItem.button.title = @"Mado";
    madoTrayStatusItem.button.toolTip = @"Mado-Tray";

    NSMenu *menu = [[NSMenu alloc] initWithTitle:@"Mado-Tray"];

    NSMenuItem *toggleItem = [[NSMenuItem alloc] initWithTitle:@"Mostrar / ocultar"
                                                        action:@selector(toggle:)
                                                 keyEquivalent:@""];
    [toggleItem setTarget:madoTrayTarget];
    [menu addItem:toggleItem];

    [menu addItem:[NSMenuItem separatorItem]];

    NSMenuItem *quitItem = [[NSMenuItem alloc] initWithTitle:@"Salir"
                                                      action:@selector(quit:)
                                               keyEquivalent:@""];
    [quitItem setTarget:madoTrayTarget];
    [menu addItem:quitItem];

    madoTrayStatusItem.menu = menu;
  });
}
