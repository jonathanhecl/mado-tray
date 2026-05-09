#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>

extern void madoTrayShow(void);
extern void madoTrayExit(void);

@interface MadoTrayMenuTarget : NSObject
- (void)showWindow:(id)sender;
- (void)exitApp:(id)sender;
@end

@implementation MadoTrayMenuTarget
- (void)showWindow:(id)sender {
  madoTrayShow();
}

- (void)exitApp:(id)sender {
  madoTrayExit();
}
@end

static NSStatusItem *madoTrayStatusItem;
static MadoTrayMenuTarget *madoTrayTarget;
static NSMenuItem *madoTrayShowItem;
static NSMenuItem *madoTrayExitItem;
static NSString *madoTrayLocale = @"en";

static NSString *MadoTrayShowLabel(void) {
  if ([madoTrayLocale isEqualToString:@"es"]) {
    return @"Mostrar";
  }
  return @"Show";
}

static NSString *MadoTrayExitLabel(void) {
  if ([madoTrayLocale isEqualToString:@"es"]) {
    return @"Salir de Mado-Tray";
  }
  return @"Exit Mado-Tray";
}

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

static void MadoTrayUpdateMenuTexts(void) {
  if (madoTrayShowItem == nil || madoTrayExitItem == nil) {
    return;
  }

  madoTrayShowItem.title = MadoTrayShowLabel();
  madoTrayExitItem.title = MadoTrayExitLabel();
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

    NSImage *image = MadoTrayIcon();
    if (image != nil) {
      madoTrayStatusItem.button.image = image;
    } else {
      madoTrayStatusItem.button.title = @"▣";
    }

    NSMenu *menu = [[NSMenu alloc] initWithTitle:@"Mado-Tray"];
    madoTrayShowItem = [[NSMenuItem alloc] initWithTitle:MadoTrayShowLabel()
                                                   action:@selector(showWindow:)
                                            keyEquivalent:@""];
    madoTrayShowItem.target = madoTrayTarget;
    [menu addItem:madoTrayShowItem];

    [menu addItem:[NSMenuItem separatorItem]];

    madoTrayExitItem = [[NSMenuItem alloc] initWithTitle:MadoTrayExitLabel()
                                                   action:@selector(exitApp:)
                                            keyEquivalent:@""];
    madoTrayExitItem.target = madoTrayTarget;
    [menu addItem:madoTrayExitItem];

    madoTrayStatusItem.menu = menu;
  });
}

void MadoTrayHide(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (madoTrayStatusItem == nil) {
      return;
    }

    [[NSStatusBar systemStatusBar] removeStatusItem:madoTrayStatusItem];
    madoTrayStatusItem = nil;
    madoTrayShowItem = nil;
    madoTrayExitItem = nil;
  });
}

void MadoTraySetLocale(char* locale) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (locale == NULL) {
      madoTrayLocale = @"en";
    } else {
      NSString *value = [NSString stringWithUTF8String:locale];
      madoTrayLocale = [value isEqualToString:@"es"] ? @"es" : @"en";
    }
    MadoTrayUpdateMenuTexts();
  });
}
