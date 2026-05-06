#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>

extern void madoTrayToggle(void);

@interface MadoTrayMenuTarget : NSObject
- (void)toggle:(id)sender;
@end

@implementation MadoTrayMenuTarget
- (void)toggle:(id)sender {
  madoTrayToggle();
}
@end

static NSStatusItem *madoTrayStatusItem;
static MadoTrayMenuTarget *madoTrayTarget;

static NSImage *MadoTrayIcon(void) {
  if (@available(macOS 11.0, *)) {
    NSImage *image = [NSImage imageWithSystemSymbolName:@"macwindow"
                               accessibilityDescription:@"Mado-Tray"];
    image.template = YES;
    return image;
  }

  return nil;
}

static void MadoTrayEnsureTarget(void) {
  if (madoTrayTarget == nil) {
    madoTrayTarget = [[MadoTrayMenuTarget alloc] init];
  }
}

void MadoTrayCreate(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    MadoTrayEnsureTarget();
  });
}

void MadoTrayShow(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    MadoTrayEnsureTarget();

    if (madoTrayStatusItem != nil) {
      return;
    }

    madoTrayStatusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    madoTrayStatusItem.button.toolTip = @"Mado-Tray";
    madoTrayStatusItem.button.target = madoTrayTarget;
    madoTrayStatusItem.button.action = @selector(toggle:);

    NSImage *image = MadoTrayIcon();
    if (image != nil) {
      madoTrayStatusItem.button.image = image;
    } else {
      madoTrayStatusItem.button.title = @"▣";
    }
  });
}

void MadoTrayHide(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (madoTrayStatusItem == nil) {
      return;
    }

    [[NSStatusBar systemStatusBar] removeStatusItem:madoTrayStatusItem];
    madoTrayStatusItem = nil;
  });
}
