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
  [NSApp activateIgnoringOtherApps:YES];
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
    NSLog(@"[Mado-Tray] creating status bar target");
    MadoTrayEnsureTarget();
  });
}

void MadoTrayShow(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    NSLog(@"[Mado-Tray] requested status bar item");
    MadoTrayEnsureTarget();

    if (madoTrayStatusItem != nil) {
      NSLog(@"[Mado-Tray] status bar item already exists");
      return;
    }

    madoTrayStatusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:32.0];
    [madoTrayStatusItem retain];
    NSLog(@"[Mado-Tray] status bar item created: %@", madoTrayStatusItem);
    madoTrayStatusItem.autosaveName = @"com.jonathanhecl.mado-tray.status-item";
    if (@available(macOS 10.12, *)) {
      madoTrayStatusItem.behavior = 0;
    }
    madoTrayStatusItem.button.toolTip = @"Mado-Tray";

    NSImage *image = MadoTrayIcon();
    if (image != nil) {
      madoTrayStatusItem.button.image = image;
    } else {
      madoTrayStatusItem.button.title = @"▣";
    }
    madoTrayStatusItem.button.enabled = YES;
    madoTrayStatusItem.button.hidden = NO;
    NSLog(@"[Mado-Tray] status bar item button: %@", madoTrayStatusItem.button);

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
    madoTrayStatusItem.length = 32.0;
    madoTrayStatusItem.visible = YES;
    madoTrayStatusItem.button.title = @"▣";
    madoTrayStatusItem.button.enabled = YES;
    madoTrayStatusItem.button.hidden = NO;
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    NSLog(@"[Mado-Tray] status bar item configured with menu");
  });
}

void MadoTrayHide(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (madoTrayStatusItem == nil) {
      return;
    }

    [[NSStatusBar systemStatusBar] removeStatusItem:madoTrayStatusItem];
    [madoTrayStatusItem release];
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
